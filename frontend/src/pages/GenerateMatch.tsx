"use client"

import { useState, useEffect } from "react"
import { Users, Settings, Play, FileText, Monitor } from "lucide-react"
import { MainLayout } from "@/components/layout"
import { 
  Tabs, 
  TabsList, 
  TabsTrigger, 
  TabsContent,
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  Button,
  Badge
} from "@/components/ui"
import { TeamBuilder } from "@/components/forms/TeamBuilder"
import { MatchSettings } from "@/components/forms/MatchSettings"
import { LogViewer } from "@/components/viewer/LogViewer"
import { 
  useAppStore, 
  useGenerationStatus, 
  useMatchConfig, 
  useTeamsManagement 
} from "@/store"
import { useMatchStream } from "@/hooks/useWebSocket"
import { useMatchStore } from "@/store/useMatchStore"
import { setToastHandler } from "@/services/api"
import type { ITeamFormData } from "@/types"
import type { ISpecificGameEvent } from "@/types/events"

export function GenerateMatch() {
  // Local state for log viewer
  const [isLogViewerFullscreen, setIsLogViewerFullscreen] = useState(false)
  const [streamingLogs, setStreamingLogs] = useState<ISpecificGameEvent[]>([])
  const [currentMatchId, setCurrentMatchId] = useState<string | null>(null)

  // Use Zustand stores
  const { 
    activeTab, 
    setActiveTab,
    generateMatch,
    canProceedToSettings,
    canGenerate,
    showToast
  } = useAppStore()
  
  const { saveGenerationResult } = useMatchStore()
  
  const { 
    status, 
    result: generationResult, 
    error: generationError,
    isGenerating 
  } = useGenerationStatus()
  
  const { config } = useMatchConfig()
  const { teams, updateTeam } = useTeamsManagement()

  // WebSocket streaming for real-time logs
  const {
    events: wsEvents,
    status: wsStatus,
    isConnected: wsConnected,
    connect: connectWebSocket,
    disconnect: disconnectWebSocket,
    pause: pauseStream,
    resume: resumeStream,
    clear: clearEvents,
    isPaused
  } = useMatchStream({
    matchId: currentMatchId || undefined,
    autoConnect: false,
    onEvent: (event) => {
      // Handle new events
      setStreamingLogs(prev => [...prev, event])
    },
    onStatusChange: (status) => {
      console.log('Stream status:', status)
    },
    onError: (error) => {
      console.error('Stream error:', error)
      showToast(`Streaming error: ${error}`, 'error')
    }
  })

  const handleTeamUpdate = (teamIndex: number, team: ITeamFormData) => {
    updateTeam(teamIndex, team)
  }

  const handleGenerate = async () => {
    try {
      // Clear previous logs and start streaming
      setStreamingLogs([])
      clearEvents()
      
      // Generate the match
      await generateMatch()
      
      // If generation started successfully, connect WebSocket for real-time streaming
      if (generationResult?.match_id) {
        setCurrentMatchId(generationResult.match_id)
        connectWebSocket()
        
        // Auto-switch to results tab to show streaming logs
        setTimeout(() => {
          setActiveTab("results")
        }, 1000)
      }
    } catch (error) {
      console.error('Generation failed:', error)
      showToast('Failed to generate match', 'error')
    }
  }

  // Update match ID when generation result changes
  useEffect(() => {
    if (generationResult?.match_id && generationResult.match_id !== currentMatchId) {
      setCurrentMatchId(generationResult.match_id)
      
      // Connect to WebSocket for streaming if we're on results tab
      if (activeTab === 'results') {
        connectWebSocket()
      }
    }
  }, [generationResult, currentMatchId, activeTab, connectWebSocket])

  // Connect WebSocket when switching to results tab with a match ID
  useEffect(() => {
    if (activeTab === 'results' && currentMatchId && !wsConnected) {
      connectWebSocket()
    }
  }, [activeTab, currentMatchId, wsConnected, connectWebSocket])

  // Cleanup WebSocket on unmount
  useEffect(() => {
    return () => {
      disconnectWebSocket()
    }
  }, [disconnectWebSocket])

  const handleToggleLogViewerFullscreen = () => {
    setIsLogViewerFullscreen(prev => !prev)
  }

  const handleDownloadLogs = () => {
    if (streamingLogs.length === 0) {
      showToast('No logs to download', 'error')
      return
    }

    const logText = streamingLogs.map(event => 
      `[${event.timestamp}] [TICK:${event.tick}] [R${event.round}] ${event.type}: ${JSON.stringify(event)}`
    ).join('\n')
    
    const blob = new Blob([logText], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `cs2-match-${currentMatchId || 'logs'}-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    showToast('Logs downloaded successfully', 'success')
  }

  // Set up toast handler for API client
  useEffect(() => {
    setToastHandler({ showToast });
  }, [showToast]);
  
  // Save generation result to match store when completed
  useEffect(() => {
    if (generationResult && status === 'completed') {
      saveGenerationResult(generationResult, {
        teams,
        map: config.map,
        format: config.format
      });
    }
  }, [generationResult, status, teams, config.map, config.format, saveGenerationResult]);

  // Combine streaming events with any static events
  const allEvents = [...streamingLogs, ...wsEvents]

  return (
    <MainLayout>
      <div className="container mx-auto px-4 py-6 max-w-6xl">
        <div className="mb-6">
          <div className="flex items-center gap-3 mb-2">
            <Play className="size-6 text-cs-orange" />
            <h1 className="text-3xl font-bold">Generate Match</h1>
          </div>
          <p className="text-muted-foreground">
            Configure teams and match settings to generate CS2 log data
          </p>
        </div>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-5 mb-6">
            <TabsTrigger value="teams" className="flex items-center gap-2">
              <Users className="size-4" />
              Teams
            </TabsTrigger>
            <TabsTrigger 
              value="settings" 
              disabled={!canProceedToSettings()}
              className="flex items-center gap-2"
            >
              <Settings className="size-4" />
              Settings
            </TabsTrigger>
            <TabsTrigger 
              value="generate" 
              disabled={!canGenerate()}
              className="flex items-center gap-2"
            >
              <Play className="size-4" />
              Generate
            </TabsTrigger>
            <TabsTrigger 
              value="results"
              disabled={!generationResult}
              className="flex items-center gap-2"
            >
              <FileText className="size-4" />
              Results
            </TabsTrigger>
            <TabsTrigger 
              value="logs"
              disabled={!generationResult || allEvents.length === 0}
              className="flex items-center gap-2"
            >
              <Monitor className="size-4" />
              Live Logs
              {wsConnected && (
                <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
              )}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="teams" className="space-y-6">
            <div className="grid gap-6 lg:grid-cols-2">
              <TeamBuilder
                team={teams[0]}
                teamIndex={0}
                title="Counter-Terrorists"
                variant="ct"
                onTeamUpdate={handleTeamUpdate}
              />
              <TeamBuilder
                team={teams[1]}
                teamIndex={1}
                title="Terrorists"
                variant="t"
                onTeamUpdate={handleTeamUpdate}
              />
            </div>
            
            <div className="flex justify-end">
              <Button
                onClick={() => setActiveTab("settings")}
                disabled={!canProceedToSettings()}
                className="min-w-32"
              >
                Next: Settings
              </Button>
            </div>
          </TabsContent>

          <TabsContent value="settings" className="space-y-6">
            <MatchSettings
              config={config}
            />
            
            <div className="flex justify-between">
              <Button
                variant="outline"
                onClick={() => setActiveTab("teams")}
              >
                Back: Teams
              </Button>
              <Button
                onClick={() => setActiveTab("generate")}
                disabled={!canGenerate()}
                className="min-w-32"
              >
                Next: Generate
              </Button>
            </div>
          </TabsContent>

          <TabsContent value="generate" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Ready to Generate</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <h4 className="font-medium mb-2">Teams</h4>
                    <div className="space-y-2">
                      {teams.map((team, index) => (
                        <div key={index} className="flex items-center gap-2">
                          <Badge variant={index === 0 ? "ct" : "t"}>
                            {index === 0 ? "CT" : "T"}
                          </Badge>
                          <span className="font-medium">{team.name}</span>
                          <span className="text-sm text-muted-foreground">
                            ({team.players.length} players)
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                  <div>
                    <h4 className="font-medium mb-2">Match Settings</h4>
                    <div className="space-y-1 text-sm">
                      <div>Map: <span className="font-medium">{config.map}</span></div>
                      <div>Format: <span className="font-medium">{config.format.toUpperCase()}</span></div>
                      <div>Overtime: <span className="font-medium">{config.overtime ? "Yes" : "No"}</span></div>
                    </div>
                  </div>
                </div>

                <div className="pt-4 border-t">
                  <Button
                    onClick={handleGenerate}
                    disabled={isGenerating || !canGenerate()}
                    className="w-full h-12 text-lg"
                  >
                    {isGenerating ? (
                      <>
                        <span className="animate-spin mr-2">⏳</span>
                        Generating Match...
                      </>
                    ) : (
                      <>
                        <Play className="mr-2 size-5" />
                        Generate CS2 Match
                      </>
                    )}
                  </Button>
                </div>
              </CardContent>
            </Card>

            <div className="flex justify-start">
              <Button
                variant="outline"
                onClick={() => setActiveTab("settings")}
                disabled={isGenerating}
              >
                Back: Settings
              </Button>
            </div>
          </TabsContent>

          <TabsContent value="results" className="space-y-6">
            {generationResult && (
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FileText className="size-5" />
                    Generation Results
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {status === "error" ? (
                    <div className="text-center py-6">
                      <div className="text-red-500 text-lg font-medium mb-2">
                        Generation Failed
                      </div>
                      <p className="text-muted-foreground mb-4">
                        {generationError || "An unknown error occurred"}
                      </p>
                      <Button
                        onClick={() => setActiveTab("generate")}
                        variant="outline"
                      >
                        Try Again
                      </Button>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      <div className="text-center py-4">
                        <div className="text-green-500 text-lg font-medium mb-2">
                          Match Generated Successfully!
                        </div>
                        <p className="text-muted-foreground mb-2">
                          Match ID: {generationResult?.match_id}
                        </p>
                        <div className="flex items-center justify-center gap-2 mb-4">
                          <Badge variant={wsConnected ? "default" : "secondary"}>
                            {wsConnected ? "🔴 Live" : "⚫ Offline"}
                          </Badge>
                          <span className="text-sm text-muted-foreground">
                            {allEvents.length} events
                          </span>
                        </div>
                        <div className="flex gap-2 justify-center">
                          {generationResult?.log_url && (
                            <Button asChild variant="outline">
                              <a href={generationResult.log_url} download>
                                Download Original
                              </a>
                            </Button>
                          )}
                          <Button 
                            onClick={handleDownloadLogs}
                            disabled={allEvents.length === 0}
                          >
                            Download Logs ({allEvents.length})
                          </Button>
                        </div>
                      </div>
                      
                      {/* Quick log preview */}
                      {allEvents.length > 0 && (
                        <div className="border rounded-lg p-4">
                          <div className="flex items-center justify-between mb-3">
                            <h4 className="font-medium">Recent Events</h4>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => setActiveTab("logs")}
                            >
                              View All Logs
                            </Button>
                          </div>
                          <div className="space-y-1 max-h-32 overflow-y-auto font-mono text-xs">
                            {allEvents.slice(-5).map((event, index) => (
                              <div key={index} className="text-gray-600">
                                <span className="text-gray-500">[{event.timestamp}]</span>
                                <span className="text-blue-400 ml-2">R{event.round}</span>
                                <span className="text-gray-400 ml-2">{event.tick}</span>
                                <span className="ml-2">{event.type}</span>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                </CardContent>
              </Card>
            )}

            <div className="flex justify-between">
              <Button
                variant="outline"
                onClick={() => {
                  setActiveTab("teams")
                  disconnectWebSocket()
                  setCurrentMatchId(null)
                  setStreamingLogs([])
                  clearEvents()
                  useAppStore.getState().clearGenerationResult()
                }}
              >
                Generate Another Match
              </Button>
              {allEvents.length > 0 && (
                <Button
                  onClick={() => setActiveTab("logs")}
                  className="flex items-center gap-2"
                >
                  <Monitor className="size-4" />
                  View Live Logs
                </Button>
              )}
            </div>
          </TabsContent>

          <TabsContent value="logs" className="space-y-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-4">
                <h2 className="text-2xl font-bold flex items-center gap-2">
                  <Monitor className="size-6" />
                  Live Match Logs
                </h2>
                <div className="flex items-center gap-2">
                  <Badge variant={wsConnected ? "default" : "secondary"}>
                    {wsConnected ? "🔴 Live" : "⚫ Offline"}
                  </Badge>
                  {currentMatchId && (
                    <Badge variant="outline">
                      ID: {currentMatchId}
                    </Badge>
                  )}
                </div>
              </div>
              
              <div className="flex items-center gap-2">
                {wsConnected && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={isPaused ? resumeStream : pauseStream}
                  >
                    {isPaused ? "Resume" : "Pause"}
                  </Button>
                )}
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    clearEvents()
                    setStreamingLogs([])
                  }}
                >
                  Clear
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleDownloadLogs}
                  disabled={allEvents.length === 0}
                >
                  Download
                </Button>
              </div>
            </div>

            <LogViewer
              events={allEvents}
              isStreaming={wsConnected && !isPaused}
              onPause={pauseStream}
              onResume={resumeStream}
              onStop={disconnectWebSocket}
              autoScroll={!isPaused}
              showControls={true}
              fullscreen={isLogViewerFullscreen}
              onToggleFullscreen={handleToggleLogViewerFullscreen}
              className="min-h-[500px]"
            />

            {allEvents.length === 0 && (
              <Card>
                <CardContent className="text-center py-12">
                  <Monitor className="size-12 mx-auto text-muted-foreground mb-4" />
                  <h3 className="text-lg font-medium mb-2">No Events Yet</h3>
                  <p className="text-muted-foreground mb-4">
                    {wsConnected 
                      ? "Waiting for match events to stream..." 
                      : "Connect to a match to see live logs"}
                  </p>
                  {!wsConnected && currentMatchId && (
                    <Button onClick={connectWebSocket}>
                      Connect to Match
                    </Button>
                  )}
                </CardContent>
              </Card>
            )}

            <div className="flex justify-start">
              <Button
                variant="outline"
                onClick={() => setActiveTab("results")}
              >
                Back to Results
              </Button>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </MainLayout>
  )
}

export default GenerateMatch