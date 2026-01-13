# Homebrew Formula Auto-Update Setup

This workflow automatically updates the Homebrew formula in `kluzzebass/homebrew-tap` when a new release is published.

## One-time Setup

You need to create a Personal Access Token (PAT) and add it as a repository secret:

### 1. Create a Personal Access Token

1. Go to https://github.com/settings/tokens/new
2. Give it a descriptive name: "Homebrew Tap Update Token"
3. Set expiration (or no expiration if you prefer)
4. Select scopes:
   - ✅ `repo` (Full control of private repositories)
5. Click "Generate token"
6. **Copy the token** (you won't be able to see it again!)

### 2. Add the token as a repository secret

1. Go to https://github.com/kluzzebass/reflag/settings/secrets/actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_TOKEN`
4. Value: Paste the token you copied
5. Click "Add secret"

## How It Works

Once set up, the workflow automatically:

1. Triggers when you publish a release (which happens automatically when you push a tag via the existing release workflow)
2. Downloads the release tarball
3. Calculates the SHA256 checksum
4. Updates `reflag.rb` in the homebrew-tap repo with the new version and SHA256
5. Commits and pushes the changes

## Usage

Just push a new tag as usual:

```bash
git tag v0.9.0
git push origin v0.9.0
```

The release workflow will create the release, then this workflow will update the Homebrew formula automatically.

## Troubleshooting

If the workflow fails:

1. Check that the `HOMEBREW_TAP_TOKEN` secret exists and is valid
2. Check the Actions tab for error logs
3. Verify the token has `repo` scope permissions
