# GitHub Pages Deployment

This branch contains the built and deployed files served via **GitHub Pages**.
Changes should only be made by the Github Action defined in the `main` branch under `/.github/workflows`.

## Folder Structure

```
.
├───dist
└───archive
    ├───v1.0.0
    |   └───dist
    └───...
```

-   `/dist`  
    Contains the files for the **current release**.

-   `/archive`  
    Contains archived builds of previous releases.

    Each version is stored under: `/archive/<VERSION>/dist`  
    where `<VERSION>` represents the release tag (e.g. `v1.0.0`, `v1.0.1`).

## Notes

-   This is an **orphan branch** intended for deployment only.
-   Source code and development files are maintained in other branches.
-   Files in this branch are generated artifacts and should not be edited manually.
