package de.psanetra.gitsemver;

import de.psanetra.gitsemver.containers.GitSemverContainer;
import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatCode;

public class NextCmdTests {

    @Test
    public void shouldIncrementVersionRelativeToSimpleTag() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.0.1");
        }

    }

    @Test
    public void shouldIncrementVersionRelativeToAnnotatedTag() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitAnnotatedTag("v1.0.0", "First Version");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("fix: Add fix");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.0.1");
        }

    }

    @Test
    public void shouldNotIncrementVersionAfterOnlyChoreCommitRelativeToSimpleTag() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("chore: Some maintenance");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.0.0");
        }

    }

    @Test
    public void shouldNotIncrementVersionAfterOnlyChoreCommitRelativeToAnnotatedTag() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitAnnotatedTag("v1.0.0", "First Version");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("chore: Some maintenance");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.0.0");
        }

    }

    @Test
    public void shouldConvertPrereleaseToRelease() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.2.3");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2");
            container.gitTag("v1.3.0-beta");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.3.0");
        }

    }

    /**
     * In this case the next command can not calculate the next version based on the commits since the latest release.
     */
    @Test
    public void shouldReturnErrorCodeIfLatestVersionNotReachableFromHEAD() throws IOException, InterruptedException {

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

            var result = container.execInContainer("git", "semver", "next");

            assertThat(result.getExitCode()).isNotEqualTo(0);
            assertThat(result.getStderr()).contains(
                "Latest tag is not on HEAD. This is necessary as the next version is calculated based on the commits since the latest version tag.");
        }

    }

    @Test
    public void shouldReturnFirstVersionOnRepoWithoutTags() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.addNewFileToGit("file.txt");
            container.gitCommit("Initial commit");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.0.0");
        }

    }

    @Test
    public void shouldReturnNextForSpecificMajorVersion() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.gitCheckoutNewBranch("v1");
            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitTag("v1.0.0");
            container.gitCheckoutNewBranch("v2");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature 2\nBREAKING CHANGE: some breaking change");
            container.gitTag("v2.0.0");
            container.gitCheckout("v1");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("fix: Fix something in v1");

            assertThat(container.exec("git", "semver", "next", "--major-version", "1")).isEqualTo("1.0.1");
        }

    }

    @Test
    public void shouldReturnLatestAfterMergeWithChoreBranchPrallelToRelease() {

        try (var container = new GitSemverContainer()) {
            container.start();

            /*
             Produce history:
             *   193c028 - Merge branch 'branch-with-chore-commit'
             |\
             * | 21c43b3 - feat: More features in master (tag: v1.0.0)
             | * 17c1f5a - chore: some maintenance (branch-with-chore-commit)
             |/
             *   1688566 - feat: First commit
             */
            container.gitCheckoutNewBranch("master");
            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: First commit");
            container.gitCheckoutNewBranch("branch-with-chore-commit");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("chore: some maintenance");
            container.gitCheckout("master");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("feat: More features in master");
            container.gitTag("v1.2.3");
            container.gitMerge("branch-with-chore-commit");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.2.3");
        }

    }

    @Test
    public void shouldReturnNewVersionAfterMergeWithFeatureBranchPrallelToRelease() {

        try (var container = new GitSemverContainer()) {
            container.start();

            /*
            Produce history:
            *   193c028 - Merge branch 'branch-with-feature-commit'
            |\
            * | 21c43b3 - feat: More features in master (tag: v1.0.0)
            | * 17c1f5a - feat: Add feature in branch (branch-with-feature-commit)
            |/
            *   1688566 - feat: Add feature
             */
            container.gitCheckoutNewBranch("master");
            container.addNewFileToGit("file.txt");
            container.gitCommit("feat: Add feature");
            container.gitCheckoutNewBranch("branch-with-feature-commit");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: Add feature in branch");
            container.gitCheckout("master");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("feat: More features in master");
            container.gitTag("v1.0.0");
            container.gitMerge("branch-with-feature-commit");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.1.0");
        }

    }

    @Test
    public void shouldPanicIfCommitIsMissingOnShallowClone() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.gitCheckoutNewBranch("master");
            container.addNewFileToGit("file.txt");
            container.gitCommit("Initial Commit");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: First Version");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("feat: Missing commit");
            container.addNewFileToGit("file4.txt");
            container.gitCommit("feat: Latest commit");
            container.exec("sh", "-c", "mv " + GitSemverContainer.WORKDIR + " /remote && mkdir " + GitSemverContainer.WORKDIR);
            container.exec("git", "clone", "--depth", "1", "file:///remote", ".");
            container.exec("git", "fetch", "--prune", "--prune-tags", "--tags");

            assertThat(container.exec("git", "log"))
                .contains("feat: Latest commit")
                .doesNotContain("feat: Missing commit")
                .doesNotContain("Initial Commit");

            assertThatCode(() -> container.exec("git", "semver", "next")).hasMessageContaining("level=fatal msg=\"object not found\"");
        }

    }

    @Test
    public void shouldReturnNewVersionOnShallowRepositoryWithAllNecessaryCommits() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.gitCheckoutNewBranch("master");
            container.addNewFileToGit("file.txt");
            container.gitCommit("Initial Commit");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: First Version");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("feat: More features");
            container.exec("sh", "-c", "mv " + GitSemverContainer.WORKDIR + " /remote && mkdir " + GitSemverContainer.WORKDIR);
            container.exec("git", "clone", "--depth", "1", "file:///remote", ".");
            container.exec("git", "fetch", "--prune", "--prune-tags", "--tags");
            container.exec("git", "fetch", "--shallow-exclude=v1.0.0");

            assertThat(container.exec("git", "log"))
                .contains("feat: More features")
                // .contains("feat: First Version")
                .doesNotContain("Initial Commit");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.1.0");
        }

    }

    @Test
    public void shouldReturnNewVersionOnShallowRepositoryWithAllNecessaryCommitsIncludingMergedBranch() {

        try (var container = new GitSemverContainer()) {
            container.start();

            container.gitCheckoutNewBranch("master");
            container.addNewFileToGit("file.txt");
            container.gitCommit("Initial Commit");
            container.gitCheckoutNewBranch("some-feature-branch");
            container.addNewFileToGit("feature-file.txt");
            container.gitCommit("feat: Feature branch commit");
            container.gitCheckout("master");
            container.addNewFileToGit("file2.txt");
            container.gitCommit("feat: First Version");
            container.gitTag("v1.0.0");
            container.addNewFileToGit("file3.txt");
            container.gitCommit("fix: Some fix in master");
            container.gitMerge("some-feature-branch");
            container.exec("sh", "-c", "mv " + GitSemverContainer.WORKDIR + " /remote && mkdir " + GitSemverContainer.WORKDIR);
            container.exec("git", "clone", "--depth", "1", "file:///remote", ".");
            container.exec("git", "fetch", "--prune", "--prune-tags", "--tags");
            container.exec("git", "fetch", "--shallow-exclude=v1.0.0");

            assertThat(container.exec("git", "log"))
                .contains("Merge branch 'some-feature-branch'")
                .contains("fix: Some fix in master")
                //  .contains("feat: First Version")
                .contains("feat: Feature branch commit")
                .doesNotContain("Initial Commit");

            assertThat(container.exec("git", "semver", "next")).isEqualTo("1.1.0");
        }

    }

}
