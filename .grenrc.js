module.exports = {
    "dataSource": "prs",
    "prefix": "",
    "onlyMilestones": false,
    "username": "weaveworks",
    "repo": "ignite",
    "groupBy": {
        "New Features": ["kind/feature"],
        "API Changes": ["kind/api-change"],
        "Enhancements": ["kind/enhancement", "area/ux"],
        "Bug Fixes": ["kind/bug"],
        "Documentation": ["kind/documentation"],
        "Testing": ["area/testing"],
        "Releasing": ["area/releasing"],
        "No category": ["closed"]
    },
    "changelogFilename": "docs/releases/next.md",
    "ignore-labels": ["kind/cleanup"],
    "template": {
        commit: ({ message, url, author, name }) => `- [${message}](${url}) - ${author ? `@${author}` : name}`,
        issue: "- {{name}} ([{{text}}]({{url}}), [@{{user_login}}]({{user_url}}))",
        label: "",
        noLabel: "closed",
        group: "\n### {{heading}}\n",
        changelogTitle: "",
        release: "## {{release}}\n\n**Released:** {{date}}\n\n{{body}}",
        releaseSeparator: "\n\n---\n\n"
    }
}