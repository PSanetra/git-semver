package de.psanetra.gitsemver.containers;

import org.testcontainers.images.builder.ImageFromDockerfile;

import java.nio.file.Paths;
import java.util.concurrent.ExecutionException;

public class BaseImage {
    public static final ImageFromDockerfile IMAGE = new ImageFromDockerfile()
        .withFileFromPath(".", Paths.get(".."));

    public static String getImageName() {
        String imageName = null;
        try {
            imageName = IMAGE.get();
        } catch (InterruptedException | ExecutionException e) {
            throw new RuntimeException(e);
        }

        return imageName;
    }
}
