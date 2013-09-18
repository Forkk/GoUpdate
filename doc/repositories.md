Outline
=======

This is a basic outline of how GoUpdate will work.


Repositories
============

In GoUpdate, the repository is where the server stores the application's update files and version information. Each application has its own repository. At the root of the repository, there is an `index.json` file with information about the repository's contents, `<version ID>.json` files for each version available, and `<channel ID>.json` for each channel.


Index
-----

The repository's index contains information about the repository's contents. It contains short, basic information about the versions and channels that the repository contains.

The index contains the following information.

+ `ApiVersion`, an integer indicating the version of the GoUpdate API that this file uses. For this version, it will be 0. 
+ `Versions`, a list of objects containing some basic information about the available versions. Each version object contains the following information.
    - `Id`, an integer used to identify this version. The version ID identifies whether this version is older or newer than another version. Any time a new version is released, this ID should be changed to a value greater than that of the previous version.
    - `Name`, a string containing the version's name. The version name is only what will be displayed to the user. It is not in any way used to determine whether or not the version is older or newer than another version.
+ `Channels`, a list of objects containing information about the available version channels. See the channels section for more information. Each object contains the following information.
    - `Id`, a string used internally to identify this channel. This isn't usually displayed to the user.
    - `Name`, the channel's display name. This is the name of the channel as it is shown to the user.


Versions
--------

The repository's versions represent different versions and iterations of the software. Each version has its own `<version ID>.json` file containing detailed information about the version and its files.

The version JSON file contains the following information.

+ `ApiVersion`, same as the `ApiVersion` field in the index.json file.
+ `Id`, integer version identifier. Must match the filename.
+ `Name`, string version display name.
+ `Files`, list of objects containing information about files that this version contains. Each object contains the following information.
    - `Path`, the path that the file will be installed to, including the filename. Relative to the program's installation directory.
    - `Sources`, an array of source objects containing information about where the file can be downloaded from in order of preference. The client should try the first source in this list that it supports and move to the next one if it fails. The update should only fail if either none of the sources for a file are supported, or all of the sources for the file fail to download successfully. See the section on source objects below.
    - `Md5`, the file's MD5 hash. Used to verify that the file was downloaded correctly. When updating, if the MD5 of an installed file on disk matches that of the same file in the version that the application is updating to, the file will not be downloaded.


### Sources

Sources in GoUpdate are a feature that allows the system to tell clients where certain files should be retreived from. There are different types of source objects with different fields, but every source object has a `source_type` field that indicates what kind of source it is. All of the supported source types are listed below.

+ `http`, an HTTP source indicates that the file should be downloaded from a certain URL over HTTP. It contains the following information.
    - `url`, the URL to download from.

+ `httpc`, an HTTP compressed source points to a compressed file that should be downloaded over HTTP and then decompressed. It contains the following information. Some clients may not support all compression types (or any compression at all), so you might want to specify additional regular `http` sources when using `httpc` sources.
    - `Url`, the URL to download from.
    - `CompressionType`, the type of compression that should be used to decompress the file. Can be any of the following: gz, bz2, xz, rar, 7z.
    - `FileInArchive`, only used with rar and 7z compression types. Indicates the name of the file inside the archive that this source refers to. Irrelevant for gz, bz2, etc. because they don't support having multiple files (unless you use tar).


Channels
--------

Channels are a feature of GoUpdate that allow an application to have multiple different release "stages" at once. For example, an application could have a "stable" channel and a "dev" channel. The dev channel contains versions of the application with experimental features that haven't been tested yet, while the stable channel contains the released, tested, and stable versions of the application. Users can then easily change between the two channels within the application depending on their preferences. Channels don't actually contain much, they simply point to a version ID that is the latest version on that channel.

Channel JSON files are stored in the root directory and contain the following information.

+ `ApiVersion`, same as the `ApiVersion` field in the index.json file.
+ `Id`, the channel's string ID. Must match the filename.
+ `Name`, the channel's display name.
+ `LatestVersion`, the latest version on this channel.


