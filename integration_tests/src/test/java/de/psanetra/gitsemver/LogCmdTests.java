package de.psanetra.gitsemver;

import de.psanetra.gitsemver.containers.GitSemverContainer;
import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatCode;

public class LogCmdTests {

    @Test
    public void shouldPrintLogWithNoTags() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log"))
                .contains("feat: Add feature")
                .contains("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintNoLogIfLatestCommitIsTagged() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");
            container.gitTag("v1.0.0");

            assertThat(container.exec("git", "semver", "log"))
                .doesNotContain("feat: Add feature")
                .doesNotContain("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogForSpecificVersion() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v0.1.0");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("fix: Add another fix");
            container.gitTag("v1.0.0");

            assertThat(container.exec("git", "semver", "log", "v1.0.0"))
                .doesNotContain("feat: Add feature")
                .contains("fix: Add fix")
                .contains("fix: Add another fix");
        }

    }

    @Test
    public void shouldPrintLogWithExcludingSimplyTaggedCommitHistory() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log"))
                .doesNotContain("feat: Add feature")
                .contains("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogWithExcludingAnnotatedTaggedCommitHistory() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitAnnotatedTag("v1.0.0", "Release 1");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log"))
                .doesNotContain("feat: Add feature")
                .contains("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogUpToSimpleTag() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log", "v1.0.0"))
                .contains("feat: Add feature")
                .doesNotContain("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogUpToAnnotatedTagged() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitAnnotatedTag("v1.0.0", "Release 1");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log", "v1.0.0"))
                .contains("feat: Add feature")
                .doesNotContain("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogWithPreReleaseCommitsInclusive() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0-alpha");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log"))
                .contains("feat: Add feature")
                .contains("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogWithPreReleaseCommitsExclusive() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0-alpha");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "log", "--exclude-pre-releases"))
                .doesNotContain("feat: Add feature")
                .contains("fix: Add fix");
        }

    }

    @Test
    public void shouldPrintLogWithLatestPrecedingVersionNotReachableFromHEAD() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Initial feature, which is contained in v1.0.0");
            container.gitCheckoutNewBranch("v1");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: v1");
            container.gitTag("v1.0.0");
            container.gitCheckout("master");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Commit which is only on v2.0.0");
            container.gitTag("v2.0.0");

            assertThat(container.exec("git", "semver", "log", "v2.0.0"))
                .doesNotContain("feat: Initial feature, which is contained in v1.0.0")
                .doesNotContain("feat: v1")
                .contains("fix: Commit which is only on v2.0.0");
        }

    }

}
