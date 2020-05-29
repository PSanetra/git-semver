package de.psanetra.gitsemver;

import de.psanetra.gitsemver.containers.GitSemverContainer;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

public class LatestCmdIncludingPreReleasesTests {

    @Test
    public void shouldReturnPrerelease() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.2.3");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2");
            container.gitTag("v1.3.0-beta");

            assertThat(container.exec("git", "semver", "latest", "--include-pre-releases")).isEqualTo("1.3.0-beta");
        }

    }

    @Test
    public void shouldReturnRegularReleaseIfLaterThanAnyPreRelease() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.2.3");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2");
            container.gitTag("v1.3.0-beta");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("feat: Add feature 3");
            container.gitTag("v1.3.0");

            assertThat(container.exec("git", "semver", "latest", "--include-pre-releases")).isEqualTo("1.3.0");
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
            container.gitTag("v1.2.3-beta");
            container.gitCheckout("master");

            assertThat(container.exec("git", "semver", "latest", "--include-pre-releases")).isEqualTo("1.2.3-beta");
        }

    }

    @Test
    public void shouldReturnEmptyVersionOnRepoWithoutTags() {

        try (var container = new GitSemverContainer()) {
            container.start();

            assertThat(container.exec("git", "semver", "latest", "--include-pre-releases")).isEqualTo("0.0.0");
        }

    }

    @Test
    public void shouldReturnLatestForSpecificMajorVersion() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.2.3-beta");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2");
            container.gitTag("v2.3.4-beta");

            assertThat(container.exec("git", "semver", "latest", "--include-pre-releases", "--major-version", "1")).isEqualTo("1.2.3-beta");
        }

    }

}
