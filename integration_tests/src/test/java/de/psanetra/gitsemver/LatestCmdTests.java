package de.psanetra.gitsemver;

import de.psanetra.gitsemver.containers.GitSemverContainer;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

public class LatestCmdTests {
    @Test
    public void shouldIgnorePrereleasesByDefault() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.2.3");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2");
            container.gitTag("v1.3.0-beta");

            assertThat(container.exec("git", "semver", "latest")).isEqualTo("1.2.3");
        }

    }

    @Test
    public void shouldReturnVersionsNotReachableFromHEAD() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.gitCheckoutNewBranch("master");
            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Master commit");
            container.gitCheckoutNewBranch("v1");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: v1 commit");
            container.gitTag("v1.2.3");
            container.gitCheckout("master");

            assertThat(container.exec("git", "semver", "latest")).isEqualTo("1.2.3");
        }

    }

    @Test
    public void shouldReturnEmptyVersionOnRepoWithoutTags() {

        try (var container = new GitSemverContainer()) {
            container.start();

            assertThat(container.exec("git", "semver", "latest")).isEqualTo("0.0.0");
        }

    }

    @Test
    public void shouldReturnLatestForSpecificMajorVersion() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.2.3");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2");
            container.gitTag("v2.3.4");

            assertThat(container.exec("git", "semver", "latest", "--major-version", "1")).isEqualTo("1.2.3");
        }

    }

}
