# Install NPM module from package.json

SWITCH on version format:
    SEMVER:
        install_from_semver
    REMOTE TAR FILE URL:
        install_from_tar_url
    LOCAL FILE PATH:
        install_from_local_path


install_from_semver:
    (local cache index lookup by semver expression)
    IF cache hit:
        (install this module from cache)
    ELSE:
        (list module versions in local cache)
        IF a version in cache fulfills semver expression:
            (add semver expression to cache index)
            (retrieve module from cache by semver expression)
            (install this module from cache)
            --END--
        ELSE
            (Generate list of module versions from remote server)
            IF a remote module version fulfills semver expression:
                (download remote module by tar url)
                (cache module by tar url)
                (cache module by semver expression)
                (install module)
            ELSE
                error!!!
                --END--


install_from_tar_url:
    (local cache index lookup by tar url)
    IF cache hit:
       (install this module from cache)
    ELSE:
        (download remote module by tar url)
        (cache module by tar url)
        (install module)


install_from_local_path
    (copy module from local path)