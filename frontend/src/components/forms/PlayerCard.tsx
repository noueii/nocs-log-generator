"use client"

import * as React from "react"
import { User, Star, Gamepad2 } from "lucide-react"
import {
  Card,
  CardContent,
  Input,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Slider,
  Badge
} from "@/components/ui"
import type { TPlayerRole } from "@/types"
import { PLAYER_ROLES } from "@/types"

interface IPlayerFormData {
  name: string
  role: TPlayerRole
  steam_id?: string
  rating: number
  country: string
}

interface PlayerCardProps {
  player: IPlayerFormData
  playerIndex: number
  teamVariant: "ct" | "t"
  onPlayerUpdate: (player: IPlayerFormData) => void
  className?: string
}

const roleDescriptions: Record<TPlayerRole, string> = {
  entry: "First to enter sites",
  awp: "Primary sniper",
  support: "Utility and support",
  lurker: "Flanker and information gatherer", 
  igl: "In-game leader and caller",
  rifler: "Core rifler and trader"
}

const roleIcons: Record<TPlayerRole, React.ReactNode> = {
  entry: <span className="text-red-500">⚡</span>,
  awp: <span className="text-blue-500">🎯</span>,
  support: <span className="text-green-500">🛡️</span>,
  lurker: <span className="text-purple-500">🥷</span>,
  igl: <span className="text-orange-500">👑</span>,
  rifler: <span className="text-yellow-500">🔫</span>
}

export function PlayerCard({
  player,
  playerIndex,
  teamVariant,
  onPlayerUpdate,
  className
}: PlayerCardProps) {
  const handleInputChange = (field: keyof IPlayerFormData, value: any) => {
    onPlayerUpdate({
      ...player,
      [field]: value
    })
  }

  const getRatingColor = (rating: number) => {
    if (rating >= 1.2) return "text-green-500"
    if (rating >= 1.0) return "text-yellow-500" 
    return "text-red-500"
  }

  const getRatingLabel = (rating: number) => {
    if (rating >= 1.3) return "Exceptional"
    if (rating >= 1.2) return "Excellent"
    if (rating >= 1.1) return "Very Good"
    if (rating >= 1.0) return "Good"
    if (rating >= 0.9) return "Average"
    if (rating >= 0.8) return "Below Average"
    return "Poor"
  }

  return (
    <Card className={`transition-all duration-200 hover:shadow-md ${className}`}>
      <CardContent className="p-4">
        <div className="space-y-4">
          {/* Header */}
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <User className="size-4 text-muted-foreground" />
              <span className="text-sm font-medium text-muted-foreground">
                Player {playerIndex + 1}
              </span>
            </div>
            <Badge variant={teamVariant} className="ml-auto">
              {teamVariant.toUpperCase()}
            </Badge>
          </div>

          {/* Player Name */}
          <div className="space-y-2">
            <Label htmlFor={`player-${playerIndex}-name`}>Player Name</Label>
            <Input
              id={`player-${playerIndex}-name`}
              placeholder="Enter player name"
              value={player.name}
              onChange={(e) => handleInputChange('name', e.target.value)}
              className="w-full"
            />
          </div>

          {/* Role Selection */}
          <div className="space-y-2">
            <Label htmlFor={`player-${playerIndex}-role`}>Role</Label>
            <Select
              value={player.role}
              onValueChange={(value: TPlayerRole) => handleInputChange('role', value)}
            >
              <SelectTrigger id={`player-${playerIndex}-role`}>
                <SelectValue placeholder="Select role">
                  {player.role && (
                    <div className="flex items-center gap-2">
                      {roleIcons[player.role]}
                      <span>{player.role.toUpperCase()}</span>
                    </div>
                  )}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {PLAYER_ROLES.map((role) => (
                  <SelectItem key={role} value={role}>
                    <div className="flex items-center gap-2">
                      {roleIcons[role]}
                      <div>
                        <div className="font-medium">{role.toUpperCase()}</div>
                        <div className="text-xs text-muted-foreground">
                          {roleDescriptions[role]}
                        </div>
                      </div>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Skill Rating */}
          <div className="space-y-3">
            <div className="flex items-center gap-2">
              <Label htmlFor={`player-${playerIndex}-rating`}>Skill Rating</Label>
              <div className="flex items-center gap-1 ml-auto">
                <Star className="size-3 text-yellow-500 fill-yellow-500" />
                <span className={`text-sm font-medium ${getRatingColor(player.rating)}`}>
                  {player.rating.toFixed(2)}
                </span>
              </div>
            </div>
            <div className="space-y-2">
              <Slider
                id={`player-${playerIndex}-rating`}
                min={0.5}
                max={2.0}
                step={0.01}
                value={[player.rating]}
                onValueChange={(values) => handleInputChange('rating', values[0])}
                className="w-full"
              />
              <div className="flex justify-between text-xs text-muted-foreground">
                <span>0.50</span>
                <span className={getRatingColor(player.rating)}>
                  {getRatingLabel(player.rating)}
                </span>
                <span>2.00</span>
              </div>
            </div>
          </div>

          {/* Steam ID (Optional) */}
          <div className="space-y-2">
            <Label htmlFor={`player-${playerIndex}-steam`} className="flex items-center gap-2">
              <Gamepad2 className="size-3" />
              Steam ID
              <span className="text-xs text-muted-foreground">(optional)</span>
            </Label>
            <Input
              id={`player-${playerIndex}-steam`}
              placeholder="STEAM_1:0:123456789"
              value={player.steam_id || ""}
              onChange={(e) => handleInputChange('steam_id', e.target.value)}
              className="w-full font-mono text-sm"
            />
          </div>

          {/* Country (simplified for MVP) */}
          <div className="space-y-2">
            <Label htmlFor={`player-${playerIndex}-country`}>Country</Label>
            <Select
              value={player.country}
              onValueChange={(value) => handleInputChange('country', value)}
            >
              <SelectTrigger id={`player-${playerIndex}-country`}>
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
        </div>
      </CardContent>
    </Card>
  )
}

export default PlayerCard