"use client"

import { Shield, Users, Tag } from "lucide-react"
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  Input,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Badge
} from "@/components/ui"
import { PlayerCard } from "./PlayerCard"
import type { TPlayerRole } from "@/types"

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

interface TeamBuilderProps {
  team: ITeamFormData
  teamIndex: number
  title: string
  variant: "ct" | "t"
  onTeamUpdate: (teamIndex: number, team: ITeamFormData) => void
  className?: string
}

const teamColors = {
  ct: "text-cs-blue border-cs-blue/20 bg-cs-blue/5",
  t: "text-cs-orange border-cs-orange/20 bg-cs-orange/5"
}

export function TeamBuilder({
  team,
  teamIndex,
  title,
  variant,
  onTeamUpdate,
  className
}: TeamBuilderProps) {
  const handleTeamFieldChange = (field: keyof ITeamFormData, value: any) => {
    onTeamUpdate(teamIndex, {
      ...team,
      [field]: value
    })
  }

  const handlePlayerUpdate = (playerIndex: number, player: IPlayerFormData) => {
    const newPlayers = [...team.players]
    newPlayers[playerIndex] = player
    onTeamUpdate(teamIndex, {
      ...team,
      players: newPlayers
    })
  }

  const generateTeamTag = (teamName: string) => {
    if (!teamName) return ""
    return teamName
      .split(' ')
      .map(word => word.charAt(0).toUpperCase())
      .join('')
      .slice(0, 4)
  }

  const getCompletionStatus = () => {
    const teamNameComplete = team.name.trim() !== ""
    const playersComplete = team.players.every(player => player.name.trim() !== "")
    
    if (teamNameComplete && playersComplete) return "complete"
    if (teamNameComplete || team.players.some(player => player.name.trim() !== "")) return "partial"
    return "empty"
  }

  const completionStatus = getCompletionStatus()
  const completionColors = {
    complete: "border-green-500/50 bg-green-500/5",
    partial: "border-yellow-500/50 bg-yellow-500/5", 
    empty: "border-muted"
  }

  return (
    <Card className={`${teamColors[variant]} ${completionColors[completionStatus]} transition-all duration-200 ${className}`}>
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <Shield className="size-5" />
            {title}
          </CardTitle>
          <div className="flex items-center gap-2">
            <Badge 
              variant={variant}
              className="text-xs"
            >
              {variant.toUpperCase()}
            </Badge>
            {completionStatus === "complete" && (
              <Badge variant="outline" className="text-green-600 border-green-500/50">
                Ready
              </Badge>
            )}
          </div>
        </div>
      </CardHeader>
      
      <CardContent className="space-y-6">
        {/* Team Information */}
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor={`team-${teamIndex}-name`} className="flex items-center gap-2">
              <Users className="size-3" />
              Team Name
            </Label>
            <Input
              id={`team-${teamIndex}-name`}
              placeholder="Enter team name"
              value={team.name}
              onChange={(e) => {
                const newName = e.target.value
                handleTeamFieldChange('name', newName)
                if (!team.tag || team.tag === generateTeamTag(team.name)) {
                  handleTeamFieldChange('tag', generateTeamTag(newName))
                }
              }}
              className="w-full"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor={`team-${teamIndex}-tag`} className="flex items-center gap-2">
              <Tag className="size-3" />
              Team Tag
            </Label>
            <Input
              id={`team-${teamIndex}-tag`}
              placeholder="TAG"
              value={team.tag}
              onChange={(e) => handleTeamFieldChange('tag', e.target.value.toUpperCase().slice(0, 4))}
              className="w-full font-mono uppercase"
              maxLength={4}
            />
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor={`team-${teamIndex}-country`}>Country</Label>
          <Select
            value={team.country}
            onValueChange={(value) => handleTeamFieldChange('country', value)}
          >
            <SelectTrigger id={`team-${teamIndex}-country`}>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="US">🇺🇸 United States</SelectItem>
              <SelectItem value="DK">🇩🇰 Denmark</SelectItem>
              <SelectItem value="SE">🇸🇪 Sweden</SelectItem>
              <SelectItem value="UA">🇺🇦 Ukraine</SelectItem>
              <SelectItem value="FR">🇫🇷 France</SelectItem>
              <SelectItem value="BR">🇧🇷 Brazil</SelectItem>
              <SelectItem value="RU">🇷🇺 Russia</SelectItem>
              <SelectItem value="FI">🇫🇮 Finland</SelectItem>
              <SelectItem value="DE">🇩🇪 Germany</SelectItem>
              <SelectItem value="PL">🇵🇱 Poland</SelectItem>
              <SelectItem value="AU">🇦🇺 Australia</SelectItem>
              <SelectItem value="CA">🇨🇦 Canada</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Players */}
        <div className="space-y-4">
          <div className="flex items-center gap-2 pt-2 border-t">
            <Users className="size-4 text-muted-foreground" />
            <h4 className="font-medium">Players ({team.players.length}/5)</h4>
          </div>
          
          <div className="grid gap-4">
            {team.players.map((player, playerIndex) => (
              <PlayerCard
                key={playerIndex}
                player={player}
                playerIndex={playerIndex}
                teamVariant={variant}
                onPlayerUpdate={(updatedPlayer) => handlePlayerUpdate(playerIndex, updatedPlayer)}
              />
            ))}
          </div>
        </div>

        {/* Quick Fill Options */}
        <div className="pt-4 border-t">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-sm font-medium text-muted-foreground">Quick Fill:</span>
          </div>
          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              onClick={() => {
                const sampleNames = variant === "ct" 
                  ? ["device", "dupreeh", "Xyp9x", "gla1ve", "Magisk"]
                  : ["s1mple", "electronic", "Perfecto", "b1t", "sdy"]
                const roles: TPlayerRole[] = ["entry", "awp", "support", "igl", "lurker"]
                
                const newPlayers = team.players.map((player, index) => ({
                  ...player,
                  name: sampleNames[index] || `Player${index + 1}`,
                  role: roles[index],
                  rating: 1.0 + (Math.random() - 0.5) * 0.4
                }))
                
                onTeamUpdate(teamIndex, { ...team, players: newPlayers })
              }}
              className="px-3 py-1 text-xs rounded-md bg-muted hover:bg-muted/80 transition-colors"
            >
              Sample Names
            </button>
            <button
              type="button"
              onClick={() => {
                const newPlayers = team.players.map((player, index) => ({
                  ...player,
                  name: player.name || `Player${index + 1}`,
                  rating: 0.8 + Math.random() * 0.6
                }))
                onTeamUpdate(teamIndex, { ...team, players: newPlayers })
              }}
              className="px-3 py-1 text-xs rounded-md bg-muted hover:bg-muted/80 transition-colors"
            >
              Random Ratings
            </button>
            <button
              type="button"
              onClick={() => {
                if (!team.name) return
                const newPlayers = team.players.map((player) => ({
                  ...player,
                  country: team.country
                }))
                onTeamUpdate(teamIndex, { ...team, players: newPlayers })
              }}
              className="px-3 py-1 text-xs rounded-md bg-muted hover:bg-muted/80 transition-colors"
            >
              Set All Countries
            </button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default TeamBuilder