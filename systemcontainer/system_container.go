package systemcontainer

import (
	"errors"
	"fmt"

	"github.com/coreos/go-systemd/dbus"
)

// Public Functions

// RunContainer runs a container with given ID
func RunSystemContainer(containerID string, containerPath string) error {
	conn, err := dbus.NewSystemConnection()
	if err != nil {
		return err
	}

	return runSysContainerInternal(containerID, containerPath, conn)
}

// Internal Functions

// Produce a per-container unique name used to identify its transient systemd service
func getUnitName(containerID string) string {
	return fmt.Sprintf("kpod-container-%s.service", containerID)
}

// Schedule a container as a transient systemd unit
func runSysContainerInternal(containerID, containerPath string, conn *dbus.Conn) error {
	if conn == nil {
		return errors.New("No connection to DBus available!")
	}

	serviceName := getUnitName(containerID)

	// Check if a unit with this name already exists
	running, err := isContainerRunning(serviceName, conn)
	if err != nil {
		return err
	} else if running {
		return fmt.Errorf("Unit with name %s already exists - container already running?", serviceName)
	}

	// The status of the systemd job, once complete, will be returned in this channel
	statusChan := make(chan string)

	// Properties for the transient unit
	// TODO do we want to set uncleanIsFailure to true? Does it even matter for a transient unit? need more info
	// TODO we're going to need to set more than just program being executed here - cgroups info, etc
	jobProperties := []dbus.Property{
		dbus.PropExecStart([]string{"/usr/bin/runc", "start", "-b", containerPath, containerID}, false),
	}

	_, err = conn.StartTransientUnit(serviceName, "fail", jobProperties, statusChan)
	if err != nil {
		return err
	}

	// Wait for the systemd job to complete. We'll get a message in the status channel
	status := <-statusChan
	if status != "done" {
		return fmt.Errorf("Error starting systemd unit for container %s", serviceName)
	}

	// TODO if the unit fails it sits around in systemd so you can get the logs with journald and such
	// This is undesirable because then we can't restart it until we 'systemctl reset-failed' manually
	// Should clean up the dead unit here

	return nil
}

// Check if the container is already running
// TODO query runc as well
// We can differentiate if the container is actually running or if we had a name collision with the systemd unit
func isContainerRunning(containerID string, conn *dbus.Conn) (bool, error) {
	if conn == nil {
		return false, errors.New("No connection to DBus available!")
	}

	serviceName := getUnitName(containerID)

	units, err := conn.ListUnits()
	if err != nil {
		return false, err
	}

	for _, unit := range units {
		if unit.Name == serviceName {
			// We found the unit, but it's not running
			// This should never happen, all our units are transient, disappear when they're done running
			// This means we have a possible name collision with a non-kpod systemd service
			if unit.ActiveState != "active" {
				return false, fmt.Errorf("Unit %s found, but not running - name collision possible", serviceName)
			}

			return true, nil
		}
	}

	return false, nil
}
