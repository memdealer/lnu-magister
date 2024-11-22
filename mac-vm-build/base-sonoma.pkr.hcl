packer {
  required_plugins {
    tart = {
      version = ">= 1.4.0"
      source  = "github.com/cirruslabs/tart"
    }
  }
}

locals {
  vm_base_name = "${var.macos_version}-${var.macos_version_tag}-base"
  vm_tag = "${var.serial_tag}"
}


variable "serial_tag" {
  type = string
}

variable "macos_version" {
  type = string
}

variable "macos_version_tag" {
  type = string
}

variable "cpu_count" {
  type = number
}

variable "memory_gb" {
  type = number
}

variable "extra_packages" {
  type = list(string)
}



variable "disk_size_gb" {
  type = number
}

variable "headless" {
  type = bool
}

variable "xcode_version" {
  type = string
}
source "tart-cli" "tart" {
  vm_base_name = "ghcr.io/cirruslabs/macos-${var.macos_version}-vanilla:${var.macos_version_tag}"
  vm_name      = local.vm_base_name
  cpu_count    = var.cpu_count
  memory_gb    = var.memory_gb
  disk_size_gb = var.disk_size_gb
  ssh_password = "admin"
  ssh_username = "admin"
  ssh_timeout  = "120s"
  headless     = var.headless
}

build {
  sources = ["source.tart-cli.tart"]

  provisioner "file" {
    source      = "data/limit.maxfiles.plist"
    destination = "~/limit.maxfiles.plist"
  }

  provisioner "file" {
    source      = "data/ntp.conf"
    destination = "/tmp/ntp.conf"
  }

  provisioner "file" {
    source      = "data/screensaver-off.sh"
    destination = "/tmp/screensaver-off.sh"
  }
  
  provisioner "file" {
    source      = "data/runner_ed25519.pub"
    destination = "/tmp/runner_ed25519.pub"
  }

    provisioner "shell" {
    inline = [
      "mkdir -p ~/.ssh",
      "chmod 755 ~/.ssh",
      "chmod 600 /tmp/runner_ed25519.pub",
      "cp -pv /tmp/runner_ed25519.pub ~/.ssh/authorized_keys",
    ]
  }

  provisioner "shell" {
    inline = [
      "echo 'Configuring maxfiles...'",
      "sudo mv ~/limit.maxfiles.plist /Library/LaunchDaemons/limit.maxfiles.plist",
      "sudo chown root:wheel /Library/LaunchDaemons/limit.maxfiles.plist",
      "sudo chmod 0644 /Library/LaunchDaemons/limit.maxfiles.plist",
      "echo 'Disabling spotlight...'",
      "sudo mdutil -a -i off",
      "sudo cp -pv /tmp/ntp.conf /etc/ntp.conf",
      "sudo ln -sf /usr/share/zoneinfo/UTC /etc/localtime"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo chsh -s /bin/bash admin",
      "sudo chsh -s /bin/bash root"
    ]
  }

  provisioner "shell" {
    inline = [
      "touch ~/.bash_profile",
      "ln -s ~/.bash_profile ~/.profile",
    ]
  }

  provisioner "shell" {
    inline = [
      "/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"",
      "echo \"export LANG=en_US.UTF-8\" >> ~/.bash_profile",
      "echo 'eval \"$(/opt/homebrew/bin/brew shellenv)\"' >> ~/.bash_profile",
      "echo \"export HOMEBREW_NO_AUTO_UPDATE=1\" >> ~/.bash_profile",
      "echo \"export HOMEBREW_NO_INSTALL_CLEANUP=1\" >> ~/.bash_profile",
      "source ~/.bash_profile",
      "brew --version",
      "brew analytics off",
      "brew update",
      "brew upgrade",
    ]
  }

  provisioner "shell" {
    inline = [
      "source ~/.bash_profile",
      "brew install ${join(" ",var.extra_packages)}"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo softwareupdate --install-rosetta --agree-to-license"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo defaults write /Library/Preferences/com.apple.loginwindow autoLoginUserScreenLocked -bool false",
      "sudo softwareupdate --schedule off",
      "sudo defaults write com.apple.SoftwareUpdate AutomaticDownload -int 0",
      "sudo defaults write com.apple.SoftwareUpdate CriticalUpdateInstall -int 0",
      "sudo defaults write com.apple.commerce AutoUpdate -bool false",
      "sudo defaults write com.apple.SoftwareUpdate AutomaticCheckEnabled -bool false",
      "sudo bash /tmp/screensaver-off.sh"
    ]
  }

  provisioner "file" {
    source      = "data/github.ssh"
    destination = "/tmp/github.ssh"
}

  provisioner "shell" {
    inline = [
      "source ~/.bash_profile",
      "[[ ! -d ~/.ssh ]] && mkdir ~/.ssh 2>/dev/null",
      "chmod 777 ~/.ssh",
      "cat /tmp/github.ssh >> ~/.ssh/known_hosts",
   ]
  }



  provisioner "shell" {
    inline = [
      "source ~/.bash_profile",
      "brew install gcc git-lfs",
      "git config --global advice.pushUpdateRejected false",
      "git config --global advice.pushNonFFCurrent false",
      "git config --global advice.pushNonFFMatching false",
      "git config --global advice.pushAlreadyExists false",
      "git config --global advice.pushFetchFirst false",
      "git config --global advice.pushNeedsForce false",
      "git config --global advice.statusHints false",
      "git config --global advice.statusUoption false",
      "git config --global advice.commitBeforeMerge false",
      "git config --global advice.resolveConflict false",
      "git config --global advice.implicitIdentity false",
      "git config --global advice.detachedHead false",
      "git config --global advice.amWorkDir false",
      "git config --global advice.rmHints false ",
    ]
  }

  provisioner "shell" {
    inline = [
      "echo 'export PATH=/usr/local/bin/:$PATH' >> ~/.bash_profile",
      "source ~/.bash_profile",
      "brew install xcodes",
      "xcodes version",
      "wget --quiet https://storage.googleapis.com/xcodes-cache/Xcode_${var.xcode_version}.xip",
      "xcodes install ${var.xcode_version} --experimental-unxip --path $PWD/Xcode_${var.xcode_version}.xip",
      "sudo rm -rf ~/.Trash/*",
      "xcodes select ${var.xcode_version}",
      "xcodebuild -downloadAllPlatforms",
      "xcodebuild -runFirstLaunch",
    ]
  }

      post-processor "manifest" {
        output = "manifest.json"
        strip_path = true
        custom_data = {
          artifact_name = local.vm_base_name
          tag = local.vm_tag
        }
    }

}
