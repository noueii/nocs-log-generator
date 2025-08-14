"use client"

import { useForm } from "react-hook-form"
import { useState } from "react"
import { Users, Settings, Play, FileText } from "lucide-react"
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
import { generateMatch } from "@/services/matchService"
import type { IMatchConfig, IGenerateResponse, TSide, TPlayerRole } from "@/types"
import { DEFAULT_MATCH_CONFIG } from "@/types"

// Use consistent form data types
interface IPlayerFormData {
  name: string
  role: TPlayerRole
  steam_id?: string
  rating: number
  country: string
}

interface ITeamFormData {
  name: string
  tag: string
  country: string
  players: IPlayerFormData[]
}

interface IMatchFormData {
  config: IMatchConfig
  teams: ITeamFormData[]
}

export function GenerateMatch() {
  const [activeTab, setActiveTab] = useState("teams")
  const [isGenerating, setIsGenerating] = useState(false)
  const [generationResult, setGenerationResult] = useState<IGenerateResponse | null>(null)
  const [teams, setTeams] = useState<ITeamFormData[]>([
    {
      name: "",
      tag: "",
      country: "US",
      players: Array(5).fill(null).map((_, i) => ({
        name: "",
        role: (i === 0 ? "entry" : i === 1 ? "awp" : i === 2 ? "support" : i === 3 ? "lurker" : "igl") as TPlayerRole,
        rating: 1.0,
        steam_id: "",
        country: "US"
      }))
    },
    {
      name: "",
      tag: "",
      country: "US", 
      players: Array(5).fill(null).map((_, i) => ({
        name: "",
        role: (i === 0 ? "entry" : i === 1 ? "awp" : i === 2 ? "support" : i === 3 ? "lurker" : "igl") as TPlayerRole,
        rating: 1.0,
        steam_id: "",
        country: "US"
      }))
    }
  ])
  
  const form = useForm<IMatchFormData>({
    defaultValues: {
      config: DEFAULT_MATCH_CONFIG,
      teams: teams
    }
  })

  const handleTeamUpdate = (teamIndex: number, team: ITeamFormData) => {
    const newTeams = [...teams]
    newTeams[teamIndex] = team
    setTeams(newTeams)
    form.setValue('teams', newTeams)
  }

  const handleConfigUpdate = (config: IMatchConfig) => {
    form.setValue('config', config)
  }

  const canProceedToSettings = () => {
    return teams.every(team => 
      team.name.trim() !== "" && 
      team.players.every(player => player.name.trim() !== "")
    )
  }

  const canGenerate = () => {
    return canProceedToSettings() && form.getValues('config')
  }

  const handleGenerate = async () => {
    if (!canGenerate()) return

    setIsGenerating(true)
    setGenerationResult(null)

    try {
      const formData = form.getValues()
      
      // Create a simplified request that matches backend expectations
      const generateRequest = {
        teams: formData.teams.map((teamData, index) => ({
          name: teamData.name,
          tag: teamData.tag,
          country: teamData.country,
          side: (index === 0 ? "CT" : "TERRORIST") as TSide,
          score: 0,
          rounds_won: 0,
          economy: {
            total_money: formData.config.start_money * 5,
            average_money: formData.config.start_money,
            equipment_value: 0,
            consecutive_losses: 0,
            loss_bonus: 1400,
            money_spent: 0,
            money_earned: 0,
            rifles: 0,
            smgs: 0,
            pistols: 5,
            snipers: 0,
            grenades: 0,
            armor: 0,
            helmets: 0,
            defuse_kits: 0
          },
          stats: {
            kills: 0,
            deaths: 0,
            assists: 0,
            score: 0,
            mvps: 0,
            adr: 0,
            first_kills: 0,
            first_deaths: 0,
            clutch_wins: 0,
            total_damage: 0
          },
          players: teamData.players.map(playerData => ({
            name: playerData.name,
            steam_id: playerData.steam_id || `STEAM_1:0:${Math.floor(Math.random() * 1000000)}`,
            role: playerData.role,
            team: teamData.name,
            side: (index === 0 ? "CT" : "TERRORIST") as TSide,
            state: {
              is_alive: true,
              health: 100,
              armor: 0,
              has_helmet: false,
              has_defuse_kit: false,
              position: { x: 0, y: 0, z: 0 },
              view_angle: { x: 0, y: 0, z: 0 },
              velocity: { x: 0, y: 0, z: 0 },
              grenades: [],
              money: formData.config.start_money,
              is_flashed: false,
              is_smoked: false,
              is_defusing: false,
              is_planting: false,
              is_reloading: false,
              has_bomb: false,
              is_last_alive: false
            },
            stats: {
              kills: 0,
              deaths: 0,
              assists: 0,
              score: 0,
              mvps: 0,
              adr: 0,
              first_kills: 0,
              first_deaths: 0,
              clutch_wins: 0,
              total_damage: 0,
              headshot_kills: 0,
              utility_damage: 0,
              enemies_flashed: 0,
              damage: 0,
              headshots: 0,
              headshot_rate: 0,
              accuracy: 0,
              trade_kills: 0,
              entry_kills: 0,
              '2k_rounds': 0,
              '3k_rounds': 0,
              '4k_rounds': 0,
              '5k_rounds': 0,
              bomb_plants: 0,
              bomb_defuses: 0,
              bomb_defuse_attempts: 0,
              hostages_rescued: 0,
              money_spent: 0,
              grenades_thrown: {},
              flash_assists: 0,
              team_kills: 0,
              team_damage: 0,
              kd_ratio: 0,
              rating: playerData.rating,
              kast: 0
            },
            economy: {
              money: formData.config.start_money,
              money_spent: 0,
              money_earned: formData.config.start_money,
              equipment_value: 0,
              purchases: [],
              eco_rounds: 0,
              force_buy_rounds: 0,
              full_buy_rounds: 0,
              economy_rating: 0
            }
          }))
        })),
        map: formData.config.map,
        format: formData.config.format,
        options: {
          seed: formData.config.seed,
          tick_rate: formData.config.tick_rate,
          overtime: formData.config.overtime,
          max_rounds: formData.config.max_rounds
        }
      } as any // Temporarily bypass strict typing for MVP

      const result = await generateMatch(generateRequest)
      setGenerationResult(result)
      setActiveTab("results")
    } catch (error) {
      console.error('Match generation failed:', error)
      setGenerationResult({
        match_id: "",
        status: "error",
        error: error instanceof Error ? error.message : "Unknown error occurred"
      })
    } finally {
      setIsGenerating(false)
    }
  }

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
          <TabsList className="grid w-full grid-cols-4 mb-6">
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
              config={form.getValues('config')}
              onConfigUpdate={handleConfigUpdate}
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
                      <div>Map: <span className="font-medium">{form.getValues('config.map')}</span></div>
                      <div>Format: <span className="font-medium">{form.getValues('config.format').toUpperCase()}</span></div>
                      <div>Overtime: <span className="font-medium">{form.getValues('config.overtime') ? "Yes" : "No"}</span></div>
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
                  {generationResult.status === "error" ? (
                    <div className="text-center py-6">
                      <div className="text-red-500 text-lg font-medium mb-2">
                        Generation Failed
                      </div>
                      <p className="text-muted-foreground mb-4">
                        {generationResult.error || "An unknown error occurred"}
                      </p>
                      <Button
                        onClick={() => setActiveTab("generate")}
                        variant="outline"
                      >
                        Try Again
                      </Button>
                    </div>
                  ) : (
                    <div className="text-center py-6">
                      <div className="text-green-500 text-lg font-medium mb-2">
                        Match Generated Successfully!
                      </div>
                      <p className="text-muted-foreground mb-4">
                        Match ID: {generationResult.match_id}
                      </p>
                      {generationResult.log_url && (
                        <Button asChild>
                          <a href={generationResult.log_url} download>
                            Download Log File
                          </a>
                        </Button>
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
                  setGenerationResult(null)
                }}
              >
                Generate Another Match
              </Button>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </MainLayout>
  )
}

export default GenerateMatch