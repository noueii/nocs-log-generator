/**
 * Error Boundary component for React error handling
 * Catches JavaScript errors anywhere in the child component tree
 */

import React, { Component, type ErrorInfo, type ReactNode } from 'react';
import { AlertTriangle, RefreshCw, Home, Bug } from 'lucide-react';
import { Button, Card, CardHeader, CardTitle, CardContent } from './ui';

/**
 * Error Boundary props
 */
interface IErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  showDetails?: boolean;
  level?: 'page' | 'component' | 'critical';
}

/**
 * Error Boundary state
 */
interface IErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
  errorId: string | null;
}

/**
 * Error Boundary class component
 */
export class ErrorBoundary extends Component<IErrorBoundaryProps, IErrorBoundaryState> {
  constructor(props: IErrorBoundaryProps) {
    super(props);
    
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    };
  }

  /**
   * Static method to update state when an error occurs
   */
  static getDerivedStateFromError(error: Error): Partial<IErrorBoundaryState> {
    // Generate unique error ID for tracking
    const errorId = `err_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    return {
      hasError: true,
      error,
      errorId,
    };
  }

  /**
   * Component did catch error lifecycle method
   */
  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Update state with error info
    this.setState({
      errorInfo,
    });

    // Log error to console in development
    if (import.meta.env.DEV || import.meta.env.VITE_ENABLE_DEBUG_LOGS === 'true') {
      console.error('🚨 Error Boundary caught an error:', error);
      console.error('Error Info:', errorInfo);
    }

    // Call custom error handler if provided
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }

    // Report error to external service in production
    this.reportError(error, errorInfo);
  }

  /**
   * Report error to external monitoring service
   */
  private reportError(error: Error, errorInfo: ErrorInfo) {
    // In a real application, you would send this to a service like Sentry, LogRocket, etc.
    const errorReport = {
      errorId: this.state.errorId,
      message: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
      timestamp: new Date().toISOString(),
      url: window.location.href,
      userAgent: navigator.userAgent,
      level: this.props.level || 'component',
    };

    // Log error report (in production, send to error reporting service)
    console.error('Error Report:', errorReport);
  }

  /**
   * Reset error state
   */
  private handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    });
  };

  /**
   * Reload page
   */
  private handleReload = () => {
    window.location.reload();
  };

  /**
   * Navigate to home page
   */
  private handleGoHome = () => {
    window.location.href = '/';
  };

  /**
   * Copy error details to clipboard
   */
  private handleCopyError = async () => {
    if (!this.state.error || !this.state.errorInfo) return;

    const errorText = `Error ID: ${this.state.errorId}
Message: ${this.state.error.message}
Stack: ${this.state.error.stack}
Component Stack: ${this.state.errorInfo.componentStack}
Timestamp: ${new Date().toISOString()}
URL: ${window.location.href}`;

    try {
      await navigator.clipboard.writeText(errorText);
      // Show success feedback (you might want to use a toast here)
      console.log('Error details copied to clipboard');
    } catch (err) {
      console.error('Failed to copy error details:', err);
    }
  };

  /**
   * Render error UI based on error level
   */
  private renderErrorUI() {
    const { level = 'component', showDetails = import.meta.env.DEV } = this.props;
    const { error, errorInfo, errorId } = this.state;

    // Critical error - full page takeover
    if (level === 'critical') {
      return (
        <div className="min-h-screen bg-background flex items-center justify-center p-4">
          <Card className="w-full max-w-2xl border-destructive">
            <CardHeader className="text-center">
              <div className="mx-auto mb-4 p-3 bg-destructive/10 rounded-full w-fit">
                <AlertTriangle className="h-8 w-8 text-destructive" />
              </div>
              <CardTitle className="text-2xl text-destructive">
                Critical Error
              </CardTitle>
              <p className="text-muted-foreground">
                The application encountered a critical error and cannot continue.
              </p>
            </CardHeader>
            <CardContent className="space-y-4">
              {showDetails && error && (
                <div className="p-4 bg-muted rounded-lg">
                  <div className="text-sm font-medium mb-2">Error Details:</div>
                  <div className="text-sm text-muted-foreground font-mono">
                    {error.message}
                  </div>
                  {errorId && (
                    <div className="text-xs text-muted-foreground mt-2">
                      Error ID: {errorId}
                    </div>
                  )}
                </div>
              )}
              <div className="flex gap-2 justify-center">
                <Button onClick={this.handleReload} variant="default">
                  <RefreshCw className="mr-2 h-4 w-4" />
                  Reload Page
                </Button>
                {showDetails && (
                  <Button onClick={this.handleCopyError} variant="outline">
                    <Bug className="mr-2 h-4 w-4" />
                    Copy Error
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      );
    }

    // Page-level error
    if (level === 'page') {
      return (
        <div className="container mx-auto px-4 py-8 max-w-4xl">
          <Card className="border-destructive">
            <CardHeader className="text-center">
              <div className="mx-auto mb-4 p-3 bg-destructive/10 rounded-full w-fit">
                <AlertTriangle className="h-6 w-6 text-destructive" />
              </div>
              <CardTitle className="text-xl text-destructive">
                Page Error
              </CardTitle>
              <p className="text-muted-foreground">
                This page encountered an error and cannot be displayed properly.
              </p>
            </CardHeader>
            <CardContent className="space-y-4">
              {showDetails && error && (
                <div className="p-3 bg-muted rounded-lg">
                  <div className="text-sm font-medium mb-1">Error:</div>
                  <div className="text-sm text-muted-foreground">
                    {error.message}
                  </div>
                  {errorId && (
                    <div className="text-xs text-muted-foreground mt-2">
                      ID: {errorId}
                    </div>
                  )}
                </div>
              )}
              <div className="flex gap-2 justify-center">
                <Button onClick={this.handleReset} variant="default">
                  <RefreshCw className="mr-2 h-4 w-4" />
                  Try Again
                </Button>
                <Button onClick={this.handleGoHome} variant="outline">
                  <Home className="mr-2 h-4 w-4" />
                  Go Home
                </Button>
                {showDetails && (
                  <Button onClick={this.handleCopyError} variant="outline" size="sm">
                    <Bug className="mr-2 h-4 w-4" />
                    Copy Error
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      );
    }

    // Component-level error - minimal UI
    return (
      <div className="p-4 border border-destructive rounded-lg bg-destructive/5">
        <div className="flex items-center gap-2 mb-2">
          <AlertTriangle className="h-4 w-4 text-destructive" />
          <span className="text-sm font-medium text-destructive">
            Component Error
          </span>
        </div>
        <p className="text-sm text-muted-foreground mb-3">
          This component failed to render properly.
        </p>
        {showDetails && error && (
          <div className="text-xs text-muted-foreground mb-3 font-mono bg-muted p-2 rounded">
            {error.message}
          </div>
        )}
        <div className="flex gap-2">
          <Button onClick={this.handleReset} variant="outline" size="sm">
            <RefreshCw className="mr-1 h-3 w-3" />
            Retry
          </Button>
          {showDetails && (
            <Button onClick={this.handleCopyError} variant="ghost" size="sm">
              <Bug className="mr-1 h-3 w-3" />
              Copy
            </Button>
          )}
        </div>
      </div>
    );
  }

  /**
   * Render method
   */
  render() {
    if (this.state.hasError) {
      // If a custom fallback is provided, use that
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // Otherwise, render built-in error UI
      return this.renderErrorUI();
    }

    // No error, render children normally
    return this.props.children;
  }
}

/**
 * Higher-order component to wrap components with ErrorBoundary
 */
export const withErrorBoundary = <T extends object>(
  Component: React.ComponentType<T>,
  errorBoundaryProps?: Omit<IErrorBoundaryProps, 'children'>
) => {
  return (props: T) => (
    <ErrorBoundary {...errorBoundaryProps}>
      <Component {...props} />
    </ErrorBoundary>
  );
};

/**
 * Hook to throw errors (for testing error boundaries)
 */
export const useErrorThrower = () => {
  return React.useCallback((error: Error | string) => {
    throw error instanceof Error ? error : new Error(error);
  }, []);
};

export default ErrorBoundary;