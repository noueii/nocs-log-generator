/**
 * Parse Demo page for uploading and parsing CS2 demo files
 */

import { useState } from 'react';
import { MainLayout } from '@/components/layout';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Upload, FileUp, AlertCircle, CheckCircle } from 'lucide-react';
import { toast } from 'sonner';

interface DemoFile {
  id: string;
  name: string;
  size: number;
  status: 'pending' | 'parsing' | 'completed' | 'error';
  progress?: number;
  error?: string;
}

export default function ParseDemo() {
  const [files, setFiles] = useState<DemoFile[]>([]);
  const [isDragging, setIsDragging] = useState(false);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const droppedFiles = Array.from(e.dataTransfer.files);
    handleFiles(droppedFiles);
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const selectedFiles = Array.from(e.target.files);
      handleFiles(selectedFiles);
    }
  };

  const handleFiles = (fileList: File[]) => {
    const demoFiles = fileList.filter(file => file.name.endsWith('.dem'));
    
    if (demoFiles.length === 0) {
      toast.error('Please select valid .dem files');
      return;
    }

    const newFiles: DemoFile[] = demoFiles.map(file => ({
      id: Math.random().toString(36).substring(7),
      name: file.name,
      size: file.size,
      status: 'pending' as const,
    }));

    setFiles(prev => [...prev, ...newFiles]);
    toast.success(`Added ${demoFiles.length} demo file(s) to queue`);
  };

  const startParsing = () => {
    if (files.length === 0) {
      toast.error('No files to parse');
      return;
    }

    toast.info('Demo parsing feature coming soon!');
    // TODO: Implement actual parsing logic
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  const getStatusIcon = (status: DemoFile['status']) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'error':
        return <AlertCircle className="h-4 w-4 text-red-500" />;
      case 'parsing':
        return <div className="h-4 w-4 border-2 border-primary border-t-transparent rounded-full animate-spin" />;
      default:
        return <FileUp className="h-4 w-4 text-gray-400" />;
    }
  };

  const getStatusBadge = (status: DemoFile['status']) => {
    switch (status) {
      case 'completed':
        return <Badge variant="success">Completed</Badge>;
      case 'error':
        return <Badge variant="destructive">Error</Badge>;
      case 'parsing':
        return <Badge variant="default">Parsing...</Badge>;
      default:
        return <Badge variant="outline">Pending</Badge>;
    }
  };

  return (
    <MainLayout>
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto space-y-8">
          {/* Header */}
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Parse Demo Files</h1>
            <p className="text-muted-foreground mt-2">
              Upload CS2 demo files (.dem) to convert them to HTTP log format
            </p>
          </div>

          {/* Upload Area */}
          <Card>
            <CardHeader>
              <CardTitle>Upload Demo Files</CardTitle>
              <CardDescription>
                Drag and drop your .dem files here or click to browse
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div
                className={`
                  relative border-2 border-dashed rounded-lg p-8 text-center transition-colors
                  ${isDragging ? 'border-primary bg-primary/5' : 'border-gray-300 hover:border-gray-400'}
                `}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
              >
                <Upload className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                <p className="text-sm font-medium mb-2">
                  Drop demo files here or click to browse
                </p>
                <p className="text-xs text-muted-foreground mb-4">
                  Only .dem files are supported
                </p>
                <Label htmlFor="file-upload" className="cursor-pointer">
                  <Button variant="outline" asChild>
                    <span>Browse Files</span>
                  </Button>
                  <Input
                    id="file-upload"
                    type="file"
                    className="hidden"
                    accept=".dem"
                    multiple
                    onChange={handleFileInput}
                  />
                </Label>
              </div>
            </CardContent>
          </Card>

          {/* File Queue */}
          {files.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>File Queue</CardTitle>
                <CardDescription>
                  {files.length} file(s) ready to parse
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {files.map(file => (
                    <div
                      key={file.id}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        {getStatusIcon(file.status)}
                        <div>
                          <p className="font-medium text-sm">{file.name}</p>
                          <p className="text-xs text-muted-foreground">
                            {formatFileSize(file.size)}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {getStatusBadge(file.status)}
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setFiles(prev => prev.filter(f => f.id !== file.id))}
                        >
                          Remove
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
                <div className="mt-4 flex gap-2">
                  <Button onClick={startParsing} disabled={files.length === 0}>
                    Start Parsing
                  </Button>
                  <Button 
                    variant="outline" 
                    onClick={() => setFiles([])}
                    disabled={files.length === 0}
                  >
                    Clear Queue
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Parsing Options */}
          <Card>
            <CardHeader>
              <CardTitle>Parsing Options</CardTitle>
              <CardDescription>
                Configure how demo files should be parsed
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="text-sm text-muted-foreground">
                <p>Demo parsing options will be available soon. Features will include:</p>
                <ul className="list-disc list-inside mt-2 space-y-1">
                  <li>Event filtering (kills, rounds, economy)</li>
                  <li>Output format selection (HTTP, JSON, CSV)</li>
                  <li>Time range selection</li>
                  <li>Player filtering</li>
                  <li>Tick rate conversion</li>
                </ul>
              </div>
            </CardContent>
          </Card>

          {/* Info */}
          <Card>
            <CardHeader>
              <CardTitle>About Demo Parsing</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground space-y-2">
              <p>
                Demo files (.dem) are recordings of CS2 matches that contain all game events,
                player positions, and server data. This tool converts them into readable HTTP
                log format for analysis.
              </p>
              <p>
                Supported demo types: POV demos, GOTV demos, and CS2 matchmaking demos.
              </p>
              <p>
                Note: This feature is currently under development and will be available soon.
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </MainLayout>
  );
}