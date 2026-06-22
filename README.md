# Mlab speedtest container application

## Clone the repository

### Clone mlab with submodules

```shell
git clone --recurse-submodules https://github.com/nokia/corteca-mlab-app.git
```

### Add submodules after cloning

In case of cloning the repository without `--recurse-submodules` you can should do:

```shell
git clone https://github.com/nokia/corteca-mlab-app.git
cd corteca-mlab-app
git submodule update --init --recursive
```

## Create the container

In order to create the container, you can use the corteca-cli. The `mlab` directory is mounted at `/app/`. All artifacts will be generated in the `./dist` directory.

### Build application with legacy containers

```shell
corteca build aarch64
corteca build armv7l
corteca build x86_64
```

### Build OCI image of container

```shell
corteca build -c 'build.options.outputType=oci' aarch64
corteca build -c 'build.options.outputType=oci' armv7l
corteca build -c 'build.options.outputType=oci' x86_64
```
