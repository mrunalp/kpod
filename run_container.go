package main

import (
	"errors"
	"fmt"

	"github.com/coreos/go-systemd/dbus"
)

// Public Functions

// RunContainer runs a container with given ID
func RunContainer(containerID string) error {
	conn, err := dbus.NewSystemConnection()
	if err != nil {
		return err
	}

	serviceName := getUnitName(containerId)

	return runContainerInternal(serviceName, conn)
}

// Internal Functions

// Produce a per-container unique name used to identify its transient systemd service
func getUnitName(containerID string) string {
	return fmt.Sprintf("kpod-container-%s.service", containerID)
}

// Schedule a container as a transient systemd unit
// TODO actually run a container and not a dummy process
func runContainerInternal(serviceName string, conn *dbus.Conn) error {
	if conn == nil {
		return errors.New("No connection to DBus available!")
	}

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
		dbus.PropExecStart([]string{"/bin/sleep", "50"}, false),
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

	return nil
}

// Check if the container is already running
// TODO query runc as well
// We can differentiate if the container is actually running or if we had a name collision with the systemd unit
func isContainerRunning(serviceName string, conn *dbus.Conn) (bool, error) {
	if conn == nil {
		return false, errors.New("No connection to DBus available!")
	}

	units, err := conn.ListUnits()
	if err != nil {
		return false, err
	}

	for _, unit := range units {
		if unit.Name == serviceName {
			return true, nil
		}
	}

	return false, nil
}
