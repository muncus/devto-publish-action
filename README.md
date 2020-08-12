# devto-publish-action

A Github Action to upload your writing to [dev.to](http://dev.to).

It supports creating new articles, or updating existing articles.

To support updating existing articles, a simple json file is required. The
suggested workflow below shows how to automatically keep this state file up
to date from your publishing process.

## Inputs

### `directory`

**Required** A directory path containing the markdown posts you would like to sync.

### `api-key`

**Required** Your Dev.to API Key. see [dev.to api docs](https://docs.dev.to/api/#section/Authentication) for instructions on creating an API Key.

### `state-file`

**Optional** the name of a file to preserve the mapping of file name to dev.to article number. The default value is `ids.json`.
In order to support updating existing articles, this file must be updated when new articles are published. See the Example Workflow below for a way to achieve this.

## Outputs

None.

## Example Action Usage

```
- name: Publish to Dev.to
    uses: muncus/devto-publish-action@release/v1
    with:
    directory: "$GITHUB_WORKSPACE/dev.to/"
    api-key: "${{ secrets.DEVTO_TOKEN }}"
    state-file: "$GITHUB_WORKSPACE/ids.json"
```
Note that this example uses a [github secret] to store the dev.to api key.

## Example Workflow

Because updating articles depends on knowing the existing Article ID number
(a requirement of the dev.to api), it is necessary to preserve changes in our
state file. To do this, I use the [peter-evans/create-pull-request action](http://github.com/peter-evans/create-pull-request).

```
    - name: Publish to Dev.to
      uses: muncus/devto-publish-action@release/v1
      with:
        directory: "$GITHUB_WORKSPACE/dev.to/"
        api-key: "${{ secrets.DEVTO_API_TOKEN }}"
        state-file: "$GITHUB_WORKSPACE/ids.json"

    - name: Create PR to update state file
      uses: peter-evans/create-pull-request@v3
      with:
        title: "[CI] Update state file with published article ids"
        reviewers: ${{ github.actor }}
        token: "${{ GITHUB_CREDENTIALS }}"
        base: master
        branch: "ci/state-file-update
```

This excerpt from my workflow creates a Pull Request whenever the state file
is updated. PRs are always based from the `master` branch, and will be
assigned to the user whose push triggered the action.

The `token` field must be filled in with a valid github auth token, whether
that's a Personal Access Token for yourself, or a dedicated automation user's Access
Token, is up to you.