package de.psanetra.gitsemver.containers;

import org.testcontainers.images.builder.ImageFromDockerfile;

import java.nio.file.Paths;

public class BaseImage {
    public static final ImageFromDockerfile IMAGE = new ImageFromDockerfile()
        .withFileFromPath(".", Paths.get(".."));

    public static String getImageName() {
        return IMAGE.get();
    }
}
