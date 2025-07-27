'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Zap, Settings, LogOut, User, Key } from 'lucide-react'

export function ChatNavbar() {
  const handleLogout = () => {
    // TODO: Implement logout
    window.location.href = '/auth/login'
  }

  const handleSettings = () => {
    // TODO: Navigate to settings
    console.log('Settings clicked')
  }

  const handleApiKeys = () => {
    // TODO: Navigate to API keys management
    window.location.href = '/auth/submit-key'
  }

  return (
    <header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50">
      <div className="flex h-16 items-center justify-between px-4">
        {/* Logo */}
        <Link href="/" className="flex items-center space-x-2">
          <Zap className="h-6 w-6 text-primary" />
          <span className="font-bold text-xl">GRYT AI</span>
        </Link>

        {/* Center - Chat Title */}
        <div className="flex-1 text-center">
          <h1 className="text-lg font-semibold text-foreground">AI Chat</h1>
        </div>

        {/* User Menu */}
        <div className="flex items-center space-x-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm" className="relative h-8 w-8 rounded-full">
                <Avatar className="h-8 w-8">
                  <AvatarImage src="" alt="User" className="" />
                  <AvatarFallback className="">
                    <User className="h-4 w-4" />
                  </AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56" align="end" forceMount>
              <div className="flex items-center justify-start gap-2 p-2">
                <div className="flex flex-col space-y-1 leading-none">
                  <p className="font-medium">User</p>
                  <p className="w-[200px] truncate text-sm text-muted-foreground">
                    user@example.com
                  </p>
                </div>
              </div>
              <DropdownMenuSeparator className="" />
              <DropdownMenuItem onClick={handleApiKeys} className="" inset={false}>
                <Key className="mr-2 h-4 w-4" />
                <span>API Keys</span>
              </DropdownMenuItem>
              <DropdownMenuItem onClick={handleSettings} className="" inset={false}>
                <Settings className="mr-2 h-4 w-4" />
                <span>Settings</span>
              </DropdownMenuItem>
              <DropdownMenuSeparator className="" />
              <DropdownMenuItem onClick={handleLogout} className="" inset={false}>
                <LogOut className="mr-2 h-4 w-4" />
                <span>Log out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}