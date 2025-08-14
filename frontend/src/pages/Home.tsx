/**
 * Home page component
 * Landing page for the CS2 Log Generator application
 */

import { Link } from 'react-router-dom';
import { 
  Play, 
  History, 
  Upload, 
  Settings, 
  FileText,
  Zap,
  Shield,
  BarChart3,
  Users,
  MapPin,
  Clock
} from 'lucide-react';
import { MainLayout } from '@/components/layout';
import { 
  Button, 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle,
  Badge 
} from '@/components/ui';
import { useMatchStatistics } from '@/store/useMatchStore';
import { useAppStore } from '@/store';

/**
 * Home page component
 */
export function Home() {
  const { showToast } = useAppStore();
  const stats = useMatchStatistics();

  const handleQuickStart = () => {
    showToast('Starting match generator...', 'info');
  };

  return (
    <MainLayout>
      <div className="container mx-auto px-4 py-8 space-y-8">
        
        {/* Hero Section */}
        <div className="text-center space-y-4">
          <div className="flex justify-center items-center gap-3 mb-4">
            <div className="p-3 bg-cs-orange/10 rounded-lg">
              <Zap className="h-8 w-8 text-cs-orange" />
            </div>
            <h1 className="text-4xl font-bold tracking-tight">
              CS2 Log Generator
            </h1>
          </div>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Generate synthetic Counter-Strike 2 match logs and parse demo files 
            for testing, analysis, and development purposes.
          </p>
          <div className="flex justify-center gap-3 mt-6">
            <Button asChild size="lg" onClick={handleQuickStart}>
              <Link to="/generate">
                <Play className="mr-2 h-5 w-5" />
                Generate Match
              </Link>
            </Button>
            <Button variant="outline" size="lg" asChild>
              <Link to="/history">
                <History className="mr-2 h-5 w-5" />
                View History
              </Link>
            </Button>
          </div>
        </div>

        {/* Stats Overview */}
        {stats.totalMatches > 0 && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Card>
              <CardContent className="p-4 text-center">
                <div className="text-2xl font-bold text-primary">
                  {stats.totalMatches}
                </div>
                <p className="text-sm text-muted-foreground">Total Matches</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <div className="text-2xl font-bold text-green-600">
                  {stats.completedMatches}
                </div>
                <p className="text-sm text-muted-foreground">Completed</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <div className="text-2xl font-bold text-cs-orange">
                  {Math.round(stats.successRate)}%
                </div>
                <p className="text-sm text-muted-foreground">Success Rate</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4 text-center">
                <div className="text-2xl font-bold text-cs-blue">
                  {Object.keys(stats.mapStats).length}
                </div>
                <p className="text-sm text-muted-foreground">Maps Used</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Main Features */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          
          {/* Generate Match */}
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <div className="p-2 bg-cs-blue/10 rounded-lg">
                  <Play className="h-5 w-5 text-cs-blue" />
                </div>
                Generate Match Logs
              </CardTitle>
              <CardDescription>
                Create synthetic CS2 match logs with customizable teams, players, and match settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <Users className="h-4 w-4 text-muted-foreground" />
                  <span>Configure teams and players</span>
                </div>
                <div className="flex items-center gap-2">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  <span>Choose from official CS2 maps</span>
                </div>
                <div className="flex items-center gap-2">
                  <Settings className="h-4 w-4 text-muted-foreground" />
                  <span>Customize match parameters</span>
                </div>
              </div>
              <div className="flex gap-2">
                <Button className="flex-1" asChild>
                  <Link to="/generate">
                    <Play className="mr-2 h-4 w-4" />
                    Start Generation
                  </Link>
                </Button>
                <Button variant="outline" asChild>
                  <Link to="/generate?tab=settings">
                    <Settings className="mr-2 h-4 w-4" />
                    Configure
                  </Link>
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Match History */}
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <div className="p-2 bg-cs-orange/10 rounded-lg">
                  <History className="h-5 w-5 text-cs-orange" />
                </div>
                Match History
              </CardTitle>
              <CardDescription>
                View, download, and manage your generated match logs and statistics
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <FileText className="h-4 w-4 text-muted-foreground" />
                  <span>Download log files</span>
                </div>
                <div className="flex items-center gap-2">
                  <BarChart3 className="h-4 w-4 text-muted-foreground" />
                  <span>View match statistics</span>
                </div>
                <div className="flex items-center gap-2">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                  <span>Track generation history</span>
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" className="flex-1" asChild>
                  <Link to="/history">
                    <History className="mr-2 h-4 w-4" />
                    View History
                  </Link>
                </Button>
                {stats.totalMatches > 0 && (
                  <Badge variant="secondary" className="px-2 py-1">
                    {stats.totalMatches}
                  </Badge>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Parse Demo */}
          <Card className="hover:shadow-lg transition-shadow opacity-75">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <div className="p-2 bg-purple-500/10 rounded-lg">
                  <Upload className="h-5 w-5 text-purple-500" />
                </div>
                Parse Demo Files
                <Badge variant="secondary" className="text-xs">
                  Coming Soon
                </Badge>
              </CardTitle>
              <CardDescription>
                Convert CS2 demo files (.dem) to HTTP log format for analysis
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <Upload className="h-4 w-4 text-muted-foreground" />
                  <span>Upload .dem files</span>
                </div>
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4 text-muted-foreground" />
                  <span>Secure file processing</span>
                </div>
                <div className="flex items-center gap-2">
                  <FileText className="h-4 w-4 text-muted-foreground" />
                  <span>Export to various formats</span>
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="secondary" disabled className="flex-1">
                  <Upload className="mr-2 h-4 w-4" />
                  Upload Demo
                </Button>
                <Button variant="outline" disabled>
                  <Settings className="mr-2 h-4 w-4" />
                  Settings
                </Button>
              </div>
            </CardContent>
          </Card>

        </div>

        {/* Recent Activity */}
        {stats.recentMatches.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="h-5 w-5" />
                Recent Matches
              </CardTitle>
              <CardDescription>
                Your latest generated matches
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {stats.recentMatches.map((match, index) => (
                  <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-muted rounded">
                        <MapPin className="h-4 w-4" />
                      </div>
                      <div>
                        <div className="font-medium">{match.id.substring(0, 12)}...</div>
                        <div className="text-sm text-muted-foreground">
                          {match.teams[0].name} vs {match.teams[1].name} on {match.map}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant={match.status === 'completed' ? 'default' : 'secondary'}>
                        {match.status}
                      </Badge>
                      <Badge variant="outline">
                        {match.format.toUpperCase()}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-4 text-center">
                <Button variant="outline" asChild>
                  <Link to="/history">
                    View All Matches
                  </Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Getting Started */}
        {stats.totalMatches === 0 && (
          <Card>
            <CardHeader>
              <CardTitle>Getting Started</CardTitle>
              <CardDescription>
                New to CS2 Log Generator? Follow these steps to create your first match.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-start gap-3">
                  <div className="flex items-center justify-center w-8 h-8 bg-primary text-primary-foreground rounded-full text-sm font-medium">
                    1
                  </div>
                  <div>
                    <h4 className="font-medium">Configure Teams</h4>
                    <p className="text-sm text-muted-foreground">
                      Set up two teams with 5 players each, including names and roles.
                    </p>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <div className="flex items-center justify-center w-8 h-8 bg-primary text-primary-foreground rounded-full text-sm font-medium">
                    2
                  </div>
                  <div>
                    <h4 className="font-medium">Choose Match Settings</h4>
                    <p className="text-sm text-muted-foreground">
                      Select map, format (MR12/MR15), and other match parameters.
                    </p>
                  </div>
                </div>
                <div className="flex items-start gap-3">
                  <div className="flex items-center justify-center w-8 h-8 bg-primary text-primary-foreground rounded-full text-sm font-medium">
                    3
                  </div>
                  <div>
                    <h4 className="font-medium">Generate & Download</h4>
                    <p className="text-sm text-muted-foreground">
                      Generate your match and download the log files for use in your projects.
                    </p>
                  </div>
                </div>
              </div>
              <div className="mt-6">
                <Button asChild>
                  <Link to="/generate">
                    <Play className="mr-2 h-4 w-4" />
                    Create Your First Match
                  </Link>
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Footer Info */}
        <div className="text-center text-sm text-muted-foreground border-t pt-6">
          <p>
            CS2 Log Generator v0.1.0 - Built for testing and development purposes
          </p>
        </div>

      </div>
    </MainLayout>
  );
}

export default Home;