package de.psanetra.gitsemver.containers;

import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.images.builder.ImageFromDockerfile;

import java.io.IOException;
import java.util.concurrent.Future;

@Slf4j
public class GitSemverContainer extends GenericContainer<GitSemverContainer> {
    public static final String WORKDIR = "/workdir";

    private static ImageFromDockerfile createImage() {
        var image = new ImageFromDockerfile()
            .withDockerfileFromBuilder(dockerfileBuilder -> {
                dockerfileBuilder
                    .from(BaseImage.getImageName())
                    .workDir(WORKDIR)
                    .entryPoint("sleep", "190");
            });

        return image;
    }

    public GitSemverContainer() {
        this(createImage());
    }

    private GitSemverContainer(@NonNull Future image) {
        super(image);
    }

    /**
     * @return stdout
     */
    public String exec(String... command) {
        ExecResult result;

        try {
            result = execInContainer(command);

            if (result.getExitCode() != 0) {
                throw new RuntimeException(
                    "Command exited with " + result.getExitCode() + ":" +
                        String.join(" ", command) +
                        "\nstdout:\n" + result.getStdout() +
                        "\nstderr:\n" + result.getStderr());
            }
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        return result.getStdout();
    }

    @Override
    public void start() {
        super.start();
        exec("git", "init");
        exec("git", "config", "user.email", "test@example.com");
        exec("git", "config", "user.name", "testuser");
    }

    public void touch(String filename) {
        exec("touch", filename);
    }

    public void gitAddAll() {
        exec("git", "add", "-A");
    }

    public void gitAdd(String filename) {
        exec("git", "add", filename);
    }

    public void addNewFileToGit(String filename) {
        touch(filename);
        gitAdd(filename);
    }

    public void gitCheckoutNewBranch(String newBranchName) {
        exec("git", "checkout", "-b", newBranchName);
    }

    public void gitCheckout(String branchName) {
        exec("git", "checkout", branchName);
    }

    public void gitCommit(String message) {
        exec("git", "commit", "-m", message);
    }

    public void gitTag(String tag) {
        exec("git", "tag", tag);
    }

    public void gitAnnotatedTag(String tag, String message) {
        exec("git", "tag", "-a", tag, "-m", message);
    }

    public void gitMerge(String branchName) {
        exec("git", "merge", branchName, "--no-edit");
    }
}
