Cloud Foundry CLI Plugin Repository (CLIPR)[![Build Status](https://travis-ci.org/cloudfoundry-incubator/cli-plugin-repo.svg?branch=master)](https://travis-ci.org/cloudfoundry-incubator/cli-plugin-repo)
=================

This is a public repository for community created CF CLI plugins. To submit your plugin
approval, please submit a pull request according to the guidelines below.

Submitting Plugins
=================
1. You need to have [git](http://git-scm.com/downloads) installed
1. Clone this repo `git clone https://github.com/cloudfoundry-incubator/cli-plugin-repo`
1. Include your plugin information in `repo-index.yml`, here is an example of a new plugin entry
  ```yaml
  - name: new_plugin
    description: new_plugin to be made available for the CF community
    version: 1.0.0
    created: 2015-1-31
    updated: 2015-1-31
    company:
    authors:
    - name: Sample-Author
      homepage: http://github.com/sample-author
      contact: contact@sample-author.io
    homepage: http://github.com/sample-author/new_plugin
    binaries:
    - platform: osx 
      url: https://github.com/sample-author/new_plugin/releases/download/v1.0.0/echo_darwin
      checksum: 2a087d5cddcfb057fbda91e611c33f46
    - platform: win64 
      url: https://github.com/sample-author/new_plugin/releases/download/v1.0.0/echo_win64.exe
      checksum: b4550d6594a3358563b9dcb81e40fd66
    - platform: linux32
      url: https://github.com/sample-author/new_plugin/releases/download/v1.0.0/echo_linux32
      checksum: f6540d6594a9684563b9lfa81e23id93
  ```
  Please make sure the spacing and colons are correct in the entry. The following descibes each field's usage.
  
  Field | Description
  ------ | ---------
  `name` | Name of your plugin, must not conflict with other existing plugins in the repo.
  `description` | Describe your plugin in a line or two. This desscription will show up when your plugin is listed on the command line
  `version` | Version number of your plugin, in [major].[minor].[build] form
  `created` | Date of first submission of the plugin, in year-month-day form
  `updated` | Date of last update of the plugin, in year-month-day form
  `company` | <b>Optional</b> field detailing company or organization that created the plugin
  `authors` | Fields to detail the authors of the plugin<br>`name`: name of author<br>`homepage`: <b>Optional</b> link to the homepage of the author<br>`contact`: <b>Optional</b> ways to contact author, email, twitter, phone etc ...
  `homepage` | Link to the homepage where the source code is hosted. Currently we only support open source plugins
  `binaries` | This section has fields detailing the various binary versions of your plugin. To reach as large an audience as possible, we encourage contributors to cross-compile their plugins on as many platforms as possible. Go provides everything you need to cross-compile for different platforms<br>`platform`: The os for this binary. Supports `osx`, `linux32`, `linux64`, `win32`, `win64`<br>`url`: Link to the binary file itself<br>`checksum`: SHA-1 of the binary file for verification

1. After making the changes, fork the repository
1. Add your fork as a remote
   ```
   cd $GOPATH/src/github.com/cloudfoundry-incubator/cli-plugin-repo
   git remote add your_name https://github.com/your_name/cli-plugin-repo
   ```
   
1. Push the changes to your fork and submit a Pull Request

 
Running your own Plugin Repo Server
=================
Included as part of this repository is the CLI Plugin Repo (CLIPR), a reference implementation of a repo server. For information on how to run CLIPR or how to write your own, [please see the CLIPR documentation here.](https://github.com/cloudfoundry-incubator/cli-plugin-repo/blob/master/docs/CLIPR.md)
