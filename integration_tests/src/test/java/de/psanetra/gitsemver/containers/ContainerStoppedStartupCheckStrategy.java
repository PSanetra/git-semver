package de.psanetra.gitsemver.containers;

import com.github.dockerjava.api.DockerClient;
import com.github.dockerjava.api.command.InspectContainerResponse;
import org.testcontainers.containers.startupcheck.StartupCheckStrategy;
import org.testcontainers.utility.DockerStatus;

import static org.testcontainers.containers.startupcheck.StartupCheckStrategy.StartupStatus.NOT_YET_KNOWN;
import static org.testcontainers.containers.startupcheck.StartupCheckStrategy.StartupStatus.SUCCESSFUL;

public class ContainerStoppedStartupCheckStrategy extends StartupCheckStrategy {

    @Override
    public StartupStatus checkStartupState(DockerClient dockerClient, String containerId) {
        InspectContainerResponse.ContainerState state = this.getCurrentState(dockerClient, containerId);
        return DockerStatus.isContainerStopped(state) ? SUCCESSFUL : NOT_YET_KNOWN;
    }

}
