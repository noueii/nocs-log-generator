"use client"

import { useState } from "react"
import { MapPin, Settings, Clock, Zap, DollarSign, Activity } from "lucide-react"
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Label,
  Slider,
  Tabs,
  TabsList,
  TabsTrigger,
  TabsContent
} from "@/components/ui"
import { useMatchConfig } from "@/store"
import type { IMatchConfig, TMapName, TMatchFormat } from "@/types"
import { MAPS } from "@/types"

interface MatchSettingsProps {
  config?: IMatchConfig
  className?: string
}

const mapDisplayNames: Record<TMapName, string> = {
  de_mirage: "Mirage",
  de_dust2: "Dust II",
  de_inferno: "Inferno", 
  de_cache: "Cache",
  de_overpass: "Overpass",
  de_train: "Train",
  de_nuke: "Nuke",
  de_cbble: "Cobblestone",
  de_vertigo: "Vertigo",
  de_ancient: "Ancient"
}

const mapDescriptions: Record<TMapName, string> = {
  de_mirage: "Balanced three-lane map with connector control",
  de_dust2: "Classic long-range dueling map",
  de_inferno: "Close-quarters apartment and site control",
  de_cache: "Open mid control with quad angles",
  de_overpass: "Vertical bathroom and connector plays",
  de_train: "Industrial site with many angles",
  de_nuke: "Vertical ramp and secret plays",
  de_cbble: "Long rotations and drop zone control",
  de_vertigo: "Unique vertical construction site",
  de_ancient: "Temple-themed with unique angles"
}

export function MatchSettings({
  config: propConfig,
  className
}: MatchSettingsProps) {
  const { config, setConfig } = useMatchConfig()
  
  // Use prop config if provided, otherwise use store config
  const activeConfig = propConfig || config
  const [selectedMap, setSelectedMap] = useState<TMapName>(activeConfig.map as TMapName)

  const handleConfigChange = (field: keyof IMatchConfig, value: any) => {
    setConfig({ [field]: value })
  }

  const handleMapSelect = (map: TMapName) => {
    setSelectedMap(map)
    handleConfigChange('map', map)
  }

  return (
    <div className={`space-y-6 ${className}`}>
      <Tabs defaultValue="basic" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="basic">Basic</TabsTrigger>
          <TabsTrigger value="economy">Economy</TabsTrigger>
          <TabsTrigger value="simulation">Simulation</TabsTrigger>
          <TabsTrigger value="advanced">Advanced</TabsTrigger>
        </TabsList>

        <TabsContent value="basic" className="space-y-6">
          {/* Map Selection */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MapPin className="size-5 text-cs-orange" />
                Map Selection
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {MAPS.map((map) => (
                  <button
                    key={map}
                    type="button"
                    onClick={() => handleMapSelect(map)}
                    className={`p-3 rounded-lg border text-left transition-all hover:shadow-md ${
                      selectedMap === map
                        ? "border-cs-orange bg-cs-orange/10"
                        : "border-muted hover:border-cs-orange/50"
                    }`}
                  >
                    <div className="font-medium">{mapDisplayNames[map]}</div>
                    <div className="text-sm text-muted-foreground mt-1">
                      {mapDescriptions[map]}
                    </div>
                    <div className="text-xs text-muted-foreground mt-2 font-mono">
                      {map}
                    </div>
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Match Format */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="size-5 text-cs-blue" />
                Match Format
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>Format</Label>
                  <Select
                    value={activeConfig.format}
                    onValueChange={(value: TMatchFormat) => handleConfigChange('format', value)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="mr12">
                        <div>
                          <div className="font-medium">MR12</div>
                          <div className="text-xs text-muted-foreground">
                            First to 13 rounds (24 rounds max)
                          </div>
                        </div>
                      </SelectItem>
                      <SelectItem value="mr15">
                        <div>
                          <div className="font-medium">MR15</div>
                          <div className="text-xs text-muted-foreground">
                            First to 16 rounds (30 rounds max)
                          </div>
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label>Overtime</Label>
                  <Select
                    value={activeConfig.overtime ? "yes" : "no"}
                    onValueChange={(value) => handleConfigChange('overtime', value === "yes")}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="no">No Overtime</SelectItem>
                      <SelectItem value="yes">Enable Overtime</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="space-y-2">
                <Label className="flex items-center gap-2">
                  <Clock className="size-3" />
                  Tick Rate: {activeConfig.tick_rate} Hz
                </Label>
                <Slider
                  value={[config.tick_rate]}
                  onValueChange={([value]) => handleConfigChange('tick_rate', value)}
                  min={64}
                  max={128}
                  step={64}
                  className="w-full"
                />
                <div className="flex justify-between text-xs text-muted-foreground">
                  <span>64 Hz (Standard)</span>
                  <span>128 Hz (High Performance)</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="economy" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <DollarSign className="size-5 text-green-500" />
                Economy Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>Start Money: ${config.start_money}</Label>
                  <Slider
                    value={[config.start_money]}
                    onValueChange={([value]) => handleConfigChange('start_money', value)}
                    min={400}
                    max={2000}
                    step={100}
                    className="w-full"
                  />
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>$400</span>
                    <span>$2000</span>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label>Max Money: ${config.max_money}</Label>
                  <Slider
                    value={[config.max_money]}
                    onValueChange={([value]) => handleConfigChange('max_money', value)}
                    min={10000}
                    max={20000}
                    step={1000}
                    className="w-full"
                  />
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>$10k</span>
                    <span>$20k</span>
                  </div>
                </div>
              </div>

              <div className="space-y-2">
                <Label>Realistic Economy</Label>
                <Select
                  value={config.realistic_economy ? "yes" : "no"}
                  onValueChange={(value) => handleConfigChange('realistic_economy', value === "yes")}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="yes">
                      <div>
                        <div className="font-medium">Realistic</div>
                        <div className="text-xs text-muted-foreground">
                          Proper save rounds and eco decisions
                        </div>
                      </div>
                    </SelectItem>
                    <SelectItem value="no">
                      <div>
                        <div className="font-medium">Simplified</div>
                        <div className="text-xs text-muted-foreground">
                          Random buy decisions
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="simulation" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Activity className="size-5 text-purple-500" />
                Simulation Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>
                  Skill Variance: {(config.skill_variance * 100).toFixed(0)}%
                </Label>
                <Slider
                  value={[config.skill_variance]}
                  onValueChange={([value]) => handleConfigChange('skill_variance', value)}
                  min={0.05}
                  max={0.5}
                  step={0.05}
                  className="w-full"
                />
                <div className="flex justify-between text-xs text-muted-foreground">
                  <span>Low Variance (Predictable)</span>
                  <span>High Variance (Chaotic)</span>
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>Rollback Events</Label>
                  <Select
                    value={config.rollback_enabled ? "yes" : "no"}
                    onValueChange={(value) => {
                      const enabled = value === "yes"
                      handleConfigChange('rollback_enabled', enabled)
                      if (!enabled) {
                        handleConfigChange('rollback_probability', 0)
                      }
                    }}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="no">Disabled</SelectItem>
                      <SelectItem value="yes">Enabled</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {config.rollback_enabled && (
                  <div className="space-y-2">
                    <Label>
                      Rollback Chance: {(config.rollback_probability * 100).toFixed(0)}%
                    </Label>
                    <Slider
                      value={[config.rollback_probability]}
                      onValueChange={([value]) => handleConfigChange('rollback_probability', value)}
                      min={0.01}
                      max={0.2}
                      step={0.01}
                      className="w-full"
                    />
                  </div>
                )}
              </div>

              <div className="space-y-4 pt-4 border-t">
                <h4 className="font-medium">Additional Events</h4>
                <div className="grid gap-3 sm:grid-cols-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={config.chat_messages}
                      onChange={(e) => handleConfigChange('chat_messages', e.target.checked)}
                      className="rounded"
                    />
                    <span className="text-sm">Chat Messages</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={config.network_issues}
                      onChange={(e) => handleConfigChange('network_issues', e.target.checked)}
                      className="rounded"
                    />
                    <span className="text-sm">Network Issues</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={config.anti_cheat_events}
                      onChange={(e) => handleConfigChange('anti_cheat_events', e.target.checked)}
                      className="rounded"
                    />
                    <span className="text-sm">Anti-Cheat Events</span>
                  </label>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="advanced" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Zap className="size-5 text-yellow-500" />
                Advanced Options
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Log Verbosity</Label>
                <Select
                  value={config.output_verbosity}
                  onValueChange={(value: 'minimal' | 'standard' | 'verbose') => 
                    handleConfigChange('output_verbosity', value)
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="minimal">
                      <div>
                        <div className="font-medium">Minimal</div>
                        <div className="text-xs text-muted-foreground">
                          Key events only
                        </div>
                      </div>
                    </SelectItem>
                    <SelectItem value="standard">
                      <div>
                        <div className="font-medium">Standard</div>
                        <div className="text-xs text-muted-foreground">
                          Normal CS2 log level
                        </div>
                      </div>
                    </SelectItem>
                    <SelectItem value="verbose">
                      <div>
                        <div className="font-medium">Verbose</div>
                        <div className="text-xs text-muted-foreground">
                          Detailed debugging info
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-4 pt-4 border-t">
                <h4 className="font-medium">Detailed Logging</h4>
                <div className="grid gap-3 sm:grid-cols-2">
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={config.include_positions}
                      onChange={(e) => handleConfigChange('include_positions', e.target.checked)}
                      className="rounded"
                    />
                    <span className="text-sm">Player Positions</span>
                  </label>
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={config.include_weapon_fire}
                      onChange={(e) => handleConfigChange('include_weapon_fire', e.target.checked)}
                      className="rounded"
                    />
                    <span className="text-sm">Weapon Fire Events</span>
                  </label>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default MatchSettings