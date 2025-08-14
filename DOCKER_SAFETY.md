# Docker Safety Guidelines

## ⚠️ IMPORTANT: Project-Specific Cleanup Only

This project's Makefile has been designed to ONLY clean up resources specific to the CS2 Log Generator project. 

### Safe Commands (Project-Specific)

These commands will ONLY affect the CS2 Log Generator containers and images:

- `make clean` - Removes ONLY this project's containers and volumes
- `make clean-images` - Removes ONLY this project's Docker images
- `make clean-all` - Removes ONLY this project's containers, volumes, and images
- `make down` - Stops and removes ONLY this project's containers

### What Gets Cleaned

When you run the cleanup commands, ONLY these resources are affected:
- Containers with names: `cs2-log-generator-backend`, `cs2-log-generator-frontend`
- Images with names: `nocs-log-generator-backend`, `nocs-log-generator-frontend`
- Volumes created by this project's docker-compose.yml
- Networks created by this project

### What is NEVER Cleaned

The cleanup commands will NEVER remove:
- Containers from other projects
- Images from other projects
- System-wide Docker resources
- Running containers (only stopped ones from this project)

## ❌ Dangerous Commands to Avoid

NEVER run these commands unless you want to affect ALL Docker resources system-wide:

```bash
# DANGEROUS - Removes ALL stopped containers from ALL projects
docker system prune

# DANGEROUS - Removes ALL unused images from ALL projects
docker image prune -a

# DANGEROUS - Removes ALL unused volumes from ALL projects
docker volume prune

# DANGEROUS - Removes everything unused system-wide
docker system prune -a --volumes
```

## ✅ Best Practices

1. Always use the project-specific Make commands for cleanup
2. Never use `docker system prune` commands directly
3. If you need to clean other projects, do so from their respective directories
4. Use `docker ps -a` to check what containers exist before cleaning
5. Use `docker images` to check what images exist before cleaning

## Recovery

If you accidentally deleted containers from other projects:

1. The containers were only stopped containers (running ones are safe)
2. You can recreate them by running `docker-compose up` in their respective project directories
3. No data in volumes should be lost unless you used `--volumes` flag
4. Images can be rebuilt or re-pulled as needed

## Project Isolation

This project is designed to be completely isolated:
- All container names are prefixed with `cs2-log-generator-`
- All image names are prefixed with `nocs-log-generator-`
- Network is named `cs2-log-generator-network`
- No shared volumes with other projects

## Support

If you need to perform system-wide Docker cleanup, use the hidden command with extreme caution:
```bash
make system-prune-careful  # Will prompt for confirmation and warn extensively
```

This command is intentionally hidden from `make help` and requires typing "yes" to confirm.