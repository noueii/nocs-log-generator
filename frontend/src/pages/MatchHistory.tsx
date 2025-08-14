"use client"

import React, { useState, useEffect } from 'react';
import { 
  History, 
  Download, 
  Trash2, 
  Eye, 
  Filter, 
  Search, 
  Calendar,
  Users,
  MapPin,
  Trophy,
  Clock,
  FileText,
  RefreshCw,
  AlertCircle
} from 'lucide-react';
import {
  Button,
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Badge,
} from '@/components/ui';

// Simple fallback components for missing UI pieces
const Dialog = ({ children }: { children: React.ReactNode }) => <div>{children}</div>;
const DialogContent = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <div className={`fixed inset-0 z-50 flex items-center justify-center bg-black/50 ${className}`}>
    <div className="bg-white rounded-lg p-6 max-w-4xl max-h-[80vh] overflow-auto">{children}</div>
  </div>
);
const DialogHeader = ({ children }: { children: React.ReactNode }) => <div className="mb-4">{children}</div>;
const DialogTitle = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <h2 className={`text-xl font-bold ${className}`}>{children}</h2>
);
const DialogTrigger = ({ children, asChild }: { children: React.ReactNode; asChild?: boolean }) => children;

const Table = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <table className={`w-full ${className}`}>{children}</table>
);
const TableHeader = ({ children }: { children: React.ReactNode }) => <thead>{children}</thead>;
const TableBody = ({ children }: { children: React.ReactNode }) => <tbody>{children}</tbody>;
const TableRow = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <tr className={`border-b ${className}`}>{children}</tr>
);
const TableHead = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <th className={`text-left p-2 font-medium ${className}`}>{children}</th>
);
const TableCell = ({ children, className }: { children: React.ReactNode; className?: string }) => (
  <td className={`p-2 ${className}`}>{children}</td>
);

// Simple pagination components
const Pagination = ({ children }: { children: React.ReactNode }) => <nav>{children}</nav>;
const PaginationContent = ({ children }: { children: React.ReactNode }) => (
  <div className="flex items-center gap-2">{children}</div>
);
const PaginationItem = ({ children }: { children: React.ReactNode }) => <div>{children}</div>;
const PaginationLink = ({ children, onClick, isActive, className }: {
  children: React.ReactNode;
  onClick?: () => void;
  isActive?: boolean;
  className?: string;
}) => (
  <button
    onClick={onClick}
    className={`px-3 py-1 rounded border ${isActive ? 'bg-blue-500 text-white' : 'bg-white'} ${className}`}
  >
    {children}
  </button>
);
const PaginationPrevious = ({ onClick, className }: { onClick?: () => void; className?: string }) => (
  <button onClick={onClick} className={`px-3 py-1 rounded border ${className}`}>
    Previous
  </button>
);
const PaginationNext = ({ onClick, className }: { onClick?: () => void; className?: string }) => (
  <button onClick={onClick} className={`px-3 py-1 rounded border ${className}`}>
    Next
  </button>
);
const PaginationEllipsis = () => <span>...</span>;
import { MainLayout } from '@/components/layout';
import { useMatchStore } from '@/store/useMatchStore';
import { matchService, type IMatch, type IMatchSearchParams } from '@/services/matchService';
import { useAppStore } from '@/store';
import { formatDistanceToNow } from 'date-fns';

/**
 * Match History page component
 */
export function MatchHistory() {
  // State
  const [matches, setMatches] = useState<IMatch[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedMatch, setSelectedMatch] = useState<IMatch | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalMatches, setTotalMatches] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [formatFilter, setFormatFilter] = useState<string>('all');
  
  // Store
  const { showToast } = useAppStore();
  const { 
    matches: storedMatches, 
    addMatch, 
    removeMatch,
    getMatchHistory 
  } = useMatchStore();

  // Load matches on mount and when filters change
  useEffect(() => {
    loadMatches();
  }, [currentPage, searchQuery, statusFilter, formatFilter]);

  /**
   * Load matches from API and store
   */
  const loadMatches = async () => {
    setIsLoading(true);
    setError(null);

    try {
      // Build search parameters
      const params: IMatchSearchParams = {
        page: currentPage,
        pageSize: 10,
        search: searchQuery || undefined,
        status: statusFilter !== 'all' ? [statusFilter] : undefined,
        format: formatFilter !== 'all' ? [formatFilter as any] : undefined,
      };

      // Try to load from API first
      const response = await matchService.getMatches(params);
      
      setMatches(response.matches);
      setTotalMatches(response.total);
      setTotalPages(Math.ceil(response.total / (params.pageSize || 10)));

    } catch (apiError) {
      console.error('Failed to load matches from API:', apiError);
      
      // Fallback to stored matches
      const localMatches = getMatchHistory();
      const filteredMatches = filterLocalMatches(localMatches);
      
      setMatches(filteredMatches.slice((currentPage - 1) * 10, currentPage * 10));
      setTotalMatches(filteredMatches.length);
      setTotalPages(Math.ceil(filteredMatches.length / 10));
      
      if (filteredMatches.length > 0) {
        showToast('Showing cached matches (offline mode)', 'info');
      } else {
        setError('No matches found and unable to connect to server');
      }
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Filter local matches based on current filters
   */
  const filterLocalMatches = (localMatches: IMatch[]) => {
    return localMatches.filter(match => {
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        if (!match.id.toLowerCase().includes(query) &&
            !match.map.toLowerCase().includes(query) &&
            !match.teams.some(team => team.name.toLowerCase().includes(query))) {
          return false;
        }
      }

      // Status filter
      if (statusFilter !== 'all' && match.status !== statusFilter) {
        return false;
      }

      // Format filter
      if (formatFilter !== 'all' && match.format !== formatFilter) {
        return false;
      }

      return true;
    });
  };

  /**
   * Handle match deletion
   */
  const handleDeleteMatch = async (matchId: string) => {
    try {
      // Try to delete from API
      await matchService.deleteMatch(matchId);
      showToast('Match deleted successfully', 'success');
    } catch (error) {
      console.error('Failed to delete match from API:', error);
      showToast('Failed to delete from server, removing locally', 'info');
    } finally {
      // Always remove from local storage
      removeMatch(matchId);
      
      // Reload matches
      await loadMatches();
    }
  };

  /**
   * Handle match download
   */
  const handleDownloadMatch = async (match: IMatch) => {
    try {
      const blob = await matchService.downloadMatchLog(match.id);
      
      // Create download link
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${match.id}_${match.map}_${match.format}.log`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      showToast('Match log downloaded', 'success');
    } catch (error) {
      console.error('Failed to download match:', error);
      showToast('Failed to download match log', 'error');
    }
  };

  /**
   * Handle search input change
   */
  const handleSearchChange = (value: string) => {
    setSearchQuery(value);
    setCurrentPage(1); // Reset to first page
  };

  /**
   * Handle filter changes
   */
  const handleFilterChange = (type: 'status' | 'format', value: string) => {
    if (type === 'status') {
      setStatusFilter(value);
    } else {
      setFormatFilter(value);
    }
    setCurrentPage(1); // Reset to first page
  };

  /**
   * Render match details modal
   */
  const renderMatchDetails = () => {
    if (!selectedMatch) return null;

    return (
      <DialogContent className="max-w-4xl max-h-[80vh] overflow-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <FileText className="h-5 w-5" />
            Match Details: {selectedMatch.id}
          </DialogTitle>
        </DialogHeader>

        <div className="grid gap-6">
          {/* Basic Info */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Match Information</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Map:</span>
                    <span>{selectedMatch.map}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Trophy className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Format:</span>
                    <Badge variant="secondary">{selectedMatch.format.toUpperCase()}</Badge>
                  </div>
                  <div className="flex items-center gap-2">
                    <Clock className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Status:</span>
                    <Badge variant={selectedMatch.status === 'completed' ? 'success' : 'secondary'}>
                      {selectedMatch.status}
                    </Badge>
                  </div>
                </div>
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Created:</span>
                    <span>{formatDistanceToNow(new Date(selectedMatch.start_time))} ago</span>
                  </div>
                  {selectedMatch.duration && (
                    <div className="flex items-center gap-2">
                      <Clock className="h-4 w-4 text-muted-foreground" />
                      <span className="font-medium">Duration:</span>
                      <span>{Math.round(selectedMatch.duration / 60)} minutes</span>
                    </div>
                  )}
                  <div className="flex items-center gap-2">
                    <FileText className="h-4 w-4 text-muted-foreground" />
                    <span className="font-medium">Events:</span>
                    <span>{selectedMatch.total_events || 0} events</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Teams */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg flex items-center gap-2">
                <Users className="h-5 w-5" />
                Teams & Scores
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid md:grid-cols-2 gap-4">
                {selectedMatch.teams.map((team, index) => (
                  <div key={index} className="p-4 border rounded-lg">
                    <div className="flex items-center justify-between mb-3">
                      <div>
                        <h4 className="font-medium">{team.name}</h4>
                        <Badge variant={team.side === 'CT' ? 'ct' : 't'} className="mt-1">
                          {team.side}
                        </Badge>
                      </div>
                      <div className="text-2xl font-bold">
                        {team.score}
                      </div>
                    </div>
                    <div className="space-y-1 text-sm">
                      {team.players.map((player, pIndex) => (
                        <div key={pIndex} className="flex justify-between">
                          <span>{player.name}</span>
                          <span className="text-muted-foreground">{player.role}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Actions */}
          <div className="flex gap-2 justify-end">
            <Button
              variant="outline"
              onClick={() => handleDownloadMatch(selectedMatch)}
            >
              <Download className="mr-2 h-4 w-4" />
              Download Log
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                handleDeleteMatch(selectedMatch.id);
                setSelectedMatch(null);
              }}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Delete Match
            </Button>
          </div>
        </div>
      </DialogContent>
    );
  };

  return (
    <MainLayout>
      <div className="container mx-auto px-4 py-6 max-w-7xl">
        {/* Header */}
        <div className="mb-6">
          <div className="flex items-center gap-3 mb-2">
            <History className="size-6 text-cs-orange" />
            <h1 className="text-3xl font-bold">Match History</h1>
          </div>
          <p className="text-muted-foreground">
            View and manage your generated CS2 matches
          </p>
        </div>

        {/* Filters */}
        <Card className="mb-6">
          <CardContent className="p-4">
            <div className="flex flex-col md:flex-row gap-4">
              <div className="flex-1">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                  <Input
                    placeholder="Search matches by ID, map, or team name..."
                    value={searchQuery}
                    onChange={(e) => handleSearchChange(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>
              <div className="flex gap-2">
                <Select value={statusFilter} onValueChange={(value) => handleFilterChange('status', value)}>
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="completed">Completed</SelectItem>
                    <SelectItem value="generating">Generating</SelectItem>
                    <SelectItem value="error">Error</SelectItem>
                  </SelectContent>
                </Select>
                <Select value={formatFilter} onValueChange={(value) => handleFilterChange('format', value)}>
                  <SelectTrigger className="w-32">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Formats</SelectItem>
                    <SelectItem value="mr12">MR12</SelectItem>
                    <SelectItem value="mr15">MR15</SelectItem>
                  </SelectContent>
                </Select>
                <Button
                  variant="outline"
                  onClick={loadMatches}
                  disabled={isLoading}
                >
                  <RefreshCw className={`mr-2 h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
                  Refresh
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Content */}
        {isLoading ? (
          <Card>
            <CardContent className="p-8 text-center">
              <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-4 text-muted-foreground" />
              <p className="text-muted-foreground">Loading matches...</p>
            </CardContent>
          </Card>
        ) : error ? (
          <Card>
            <CardContent className="p-8 text-center">
              <AlertCircle className="h-8 w-8 mx-auto mb-4 text-destructive" />
              <p className="text-destructive mb-4">{error}</p>
              <Button onClick={loadMatches} variant="outline">
                <RefreshCw className="mr-2 h-4 w-4" />
                Try Again
              </Button>
            </CardContent>
          </Card>
        ) : matches.length === 0 ? (
          <Card>
            <CardContent className="p-8 text-center">
              <History className="h-8 w-8 mx-auto mb-4 text-muted-foreground" />
              <p className="text-muted-foreground mb-4">No matches found</p>
              <Button onClick={() => window.location.href = '/generate'}>
                Generate Your First Match
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-6">
            {/* Matches Table */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span>Matches ({totalMatches})</span>
                  <Badge variant="secondary">{matches.length} showing</Badge>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Match ID</TableHead>
                      <TableHead>Teams</TableHead>
                      <TableHead>Map</TableHead>
                      <TableHead>Format</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {matches.map((match) => (
                      <TableRow key={match.id}>
                        <TableCell className="font-mono text-sm">
                          {match.id.substring(0, 8)}...
                        </TableCell>
                        <TableCell>
                          <div className="space-y-1">
                            {match.teams.map((team, index) => (
                              <div key={index} className="flex items-center gap-2">
                                <Badge variant={team.side === 'CT' ? 'ct' : 't'} className="text-xs">
                                  {team.side}
                                </Badge>
                                <span className="text-sm">{team.name}</span>
                                <span className="text-sm font-medium">{team.score}</span>
                              </div>
                            ))}
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <MapPin className="h-4 w-4 text-muted-foreground" />
                            {match.map}
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="secondary">
                            {match.format.toUpperCase()}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge variant={match.status === 'completed' ? 'success' : 
                                        match.status === 'error' ? 'destructive' : 'secondary'}>
                            {match.status}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                          {formatDistanceToNow(new Date(match.start_time))} ago
                        </TableCell>
                        <TableCell>
                          <div className="flex gap-1">
                            <Dialog>
                              <DialogTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => setSelectedMatch(match)}
                                >
                                  <Eye className="h-4 w-4" />
                                </Button>
                              </DialogTrigger>
                              {renderMatchDetails()}
                            </Dialog>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleDownloadMatch(match)}
                            >
                              <Download className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleDeleteMatch(match.id)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center">
                <Pagination>
                  <PaginationContent>
                    <PaginationItem>
                      <PaginationPrevious
                        onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                        className={currentPage === 1 ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                      />
                    </PaginationItem>
                    
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      const page = i + 1;
                      return (
                        <PaginationItem key={page}>
                          <PaginationLink
                            onClick={() => setCurrentPage(page)}
                            isActive={currentPage === page}
                            className="cursor-pointer"
                          >
                            {page}
                          </PaginationLink>
                        </PaginationItem>
                      );
                    })}
                    
                    {totalPages > 5 && (
                      <>
                        <PaginationItem>
                          <PaginationEllipsis />
                        </PaginationItem>
                        <PaginationItem>
                          <PaginationLink
                            onClick={() => setCurrentPage(totalPages)}
                            isActive={currentPage === totalPages}
                            className="cursor-pointer"
                          >
                            {totalPages}
                          </PaginationLink>
                        </PaginationItem>
                      </>
                    )}
                    
                    <PaginationItem>
                      <PaginationNext
                        onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                        className={currentPage === totalPages ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                      />
                    </PaginationItem>
                  </PaginationContent>
                </Pagination>
              </div>
            )}
          </div>
        )}
      </div>
    </MainLayout>
  );
}

export default MatchHistory;