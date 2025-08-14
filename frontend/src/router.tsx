/**
 * Application router configuration
 * Sets up routes for the CS2 Log Generator application
 */

import { createBrowserRouter, Navigate } from 'react-router-dom';
import { ErrorBoundary } from '@/components/ErrorBoundary';

// Lazy load pages for better performance
import { lazy, Suspense } from 'react';
import { Card, CardContent } from '@/components/ui';
import { MainLayout } from '@/components/layout';

const HomePage = lazy(() => import('@/pages/Home'));
const GenerateMatch = lazy(() => import('@/pages/GenerateMatch'));
const MatchHistory = lazy(() => import('@/pages/MatchHistory'));
const ParseDemo = lazy(() => import('@/pages/ParseDemo').catch(() => ({ default: () => <div>Parse Demo page coming soon...</div> })));

/**
 * Loading component for suspense
 */
const PageLoader = ({ message = 'Loading...' }: { message?: string }) => (
  <MainLayout>
    <div className="container mx-auto px-4 py-8">
      <Card>
        <CardContent className="p-8 text-center">
          <div className="animate-spin mx-auto mb-4 h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
          <p className="text-muted-foreground">{message}</p>
        </CardContent>
      </Card>
    </div>
  </MainLayout>
);

/**
 * Error page component
 */
const ErrorPage = () => (
  <MainLayout>
    <div className="container mx-auto px-4 py-8">
      <Card>
        <CardContent className="p-8 text-center">
          <h1 className="text-2xl font-bold text-destructive mb-4">Page Not Found</h1>
          <p className="text-muted-foreground mb-4">
            The page you're looking for doesn't exist.
          </p>
          <a href="/" className="text-primary hover:underline">
            Go back to home
          </a>
        </CardContent>
      </Card>
    </div>
  </MainLayout>
);

/**
 * Wrapper component with error boundary and suspense
 */
const PageWrapper = ({ 
  children, 
  errorLevel = 'page' 
}: { 
  children: React.ReactNode;
  errorLevel?: 'page' | 'component' | 'critical';
}) => (
  <ErrorBoundary level={errorLevel} showDetails={import.meta.env.DEV}>
    <Suspense fallback={<PageLoader />}>
      {children}
    </Suspense>
  </ErrorBoundary>
);

/**
 * Application router
 */
export const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <PageWrapper>
        <HomePage />
      </PageWrapper>
    ),
    errorElement: <ErrorPage />,
  },
  {
    path: '/generate',
    element: (
      <PageWrapper>
        <GenerateMatch />
      </PageWrapper>
    ),
    errorElement: <ErrorPage />,
  },
  {
    path: '/history',
    element: (
      <PageWrapper>
        <MatchHistory />
      </PageWrapper>
    ),
    errorElement: <ErrorPage />,
  },
  {
    path: '/parse',
    element: (
      <PageWrapper>
        <ParseDemo />
      </PageWrapper>
    ),
    errorElement: <ErrorPage />,
  },
  {
    path: '/match/:id',
    element: (
      <PageWrapper>
        <MatchHistory />
      </PageWrapper>
    ),
    errorElement: <ErrorPage />,
  },
  // Redirect common paths
  {
    path: '/matches',
    element: <Navigate to="/history" replace />,
  },
  {
    path: '/generator',
    element: <Navigate to="/generate" replace />,
  },
  // Catch-all route
  {
    path: '*',
    element: <ErrorPage />,
  },
]);

export default router;