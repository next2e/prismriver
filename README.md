# Prismriver
> A revamped communal music player.

Prismriver is a rewrite of the [2E-Music](https://github.com/next2e/2E-Music) 
project, which is  currently used to play music in one of the lounges at my 
current college dorm.

Compared to the original project, there are a number of changes:
- Written in Go instead of Python
- Uses WebSockets instead of setInterval AJAX queries to obtain player state
- Integrates directly with libvlc
- Downloads songs much more quickly and uses 
[Opus](https://en.wikipedia.org/wiki/Opus_(audio_format)) for compression
- Deployable as a singular binary and as a Docker container (provided that 
ALSA is used on the host)
- Other additional features compared to the original project

> **Contents**
> - [Usage](#usage)
> - [API](#api)
> - [Installation](#installation)
>     - [Docker](#docker)
>     - [Compiling](#compiling)


# Usage
The following environment variables can be used to configure the behavior of 
the container:

| Variable Name | Description | Default Value |
| --- | --- | --- |
| PRISMRIVER_DATADIR | The directory used to store data. | /var/lib/prismriver|
| PRISMRIVER_DB_HOST | The hostname of the Postgres server. | localhost |
| PRISMRIVER_DB_NAME | The name of the Postgres database. | prismriver |
| PRISMRIVER_DB_PASSWORD | The password to authenticate with the database. | prismriver |
| PRISMRIVER_DB_PORT | The port of the Postgres server. | 5432 |
| PRISMRIVER_DB_USER | The username to authenticate with the database. | prismriver
| PRISMRIVER_VERBOSITY | The logging level of the server. | info |
These can either be specified in your command when running the server, as flags
when launching a Docker container, or through other means.

The web UI can be accessed at port 80 on whatever address is serving 
Prismriver. This UI is the primary way of interacting with the server, whether
it be playing music or managing the current play queue.

# API
API documentation is currently in the works.

# Installation
Prismriver can be installed in two ways: as a Docker container and as a binary.

Regardless of which installation method you select, there are a few 
commonalities between both in terms of setup:
- A working Postgres server must be available for Prismriver to store data
about downloaded tracks.

Once this has been set up, go to the corresponding section depending on if you
want to use Docker or manual compilation to install Prismriver.
## Docker
The latest builds of Prismriver can be obtained from Quay.io:
```bash
docker pull quay.io/ttpcodes/prismriver
```
These builds come with most of the required dependencies pre-installed, so all
that is needed is to configure various options for the Docker container to work
properly.

By default, the container exposes the HTTP frontend on port 80. This port can
be published to make the frontend publicly accessible.

For audio to work properly, the Docker host must have a working ALSA
installation with a free device available for the container to use. To pass
through the devices on the host, the `/dev/snd/` directory must be mounted in
the container. **The Docker container will not work with PulseAudio.**

Data in the container is stored at `/var/lib/prismriver` by default and can be
mounted as a Docker volume to preserve data. Alternatively, this container can
also be configured using the environment variables described under 
[Usage](#usage).

An example launch command for the Docker container is shown here:
```bash
docker run -it -v /some/host/directory:/var/lib/prismriver -p 80:80 \
--device /dev/snd --name prismriver  quay.io/ttpcodes/prismriver
```

For general information on configuring Prismriver, see [Usage](#usage).
## Compiling
Alternatively, you can compile the binary for prismriver yourself. To compile
the binary, you will need the following dependencies:
- [Working Go environment](https://golang.org/) (for compiling the backend 
code)
    - [dep](https://github.com/golang/dep) (for retrieving backend 
    dependencies)
    - [statik](https://github.com/rakyll/statik) (for bundling web assets into 
    Go)
- [Working Node.js installation](https://nodejs.org/) (for compiling frontend 
assets)
    - [yarn](https://yarnpkg.com/) (for retrieving frontend dependencies)
- git (for cloning this repo)
- make (for running the Makefile commands)
- libvlc headers (for compiling audio support)

Begin by cloning the repository:
```bash
git clone https://gitlab.com/ttpcodes/prismriver
cd prismriver/
```
This project contains a Makefile which will handle building of everything for
you, including frontend assets, backend assets, and even dependencies. These
are all of the Makefile commands:

| Command | Description | Usage |
| --- | --- | --- |
| all | Calls deps and build. | make |
| build | Builds the final binary. Calls frontend. | make build |
| deps | Installs all dep and yarn dependencies. | make deps
| frontend | Builds and bundles the frontend. | make frontend
| install | Installs the binary to /usr/local/bin. | sudo make install
| run | Runs the binary (for development purposes only). Calls build. | make run |
To build and install Prismriver, run these commands:
```bash
# Install all dep and yarn dependencies and build the binary:
make
# Install the binary to /usr/local/bin:
sudo make install
```
Before running the binary, there are a few dependencies you will need to
install manually:
- ffmpeg (for transcoding downloaded audio files)
- libopus (for transcoding downloaded audio files)
- vlc (for audio playback)
- some kind of working sound system, such as ALSA or PulseAudio

At this point, you can execute the binary to run the server.

For general information on configuring Prismriver, see [Usage](#usage).