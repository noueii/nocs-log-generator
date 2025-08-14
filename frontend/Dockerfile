# Multi-stage Dockerfile for React frontend
# Stage 1: Development stage with hot reload
FROM node:20-alpine AS dev
LABEL stage=development

# Set working directory
WORKDIR /app

# Install dependencies first for better caching
COPY package*.json ./

# Use npm for consistency (can switch to pnpm if preferred)
RUN npm install

# Copy source code
COPY . .

# Expose Vite dev server port
EXPOSE 5173

# Health check for development
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:5173 || exit 1

# Start development server with host binding for Docker
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Stage 2: Builder stage for production build
FROM node:20-alpine AS builder
LABEL stage=build

# Set working directory
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Stage 3: Production stage with Nginx
FROM nginx:1.25-alpine AS production
LABEL stage=production

# Remove default Nginx configuration
RUN rm -rf /usr/share/nginx/html/*

# Copy built application from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Create non-root user for security
RUN adduser -D -s /bin/sh -g 'nginx' nginxuser

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80 || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]