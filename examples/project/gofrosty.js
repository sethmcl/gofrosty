module.exports = {
    parseDependency: function (dependency, module) {
        if (module.Name === 'lib') {
            module.SetDownloadURLs(['whats up']);
        }
    }
};
