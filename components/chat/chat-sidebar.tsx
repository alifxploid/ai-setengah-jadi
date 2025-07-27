'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Plus, MessageSquare, MoreVertical, Trash2, Edit, X } from 'lucide-react'
import { cn } from '@/lib/utils'

interface ChatSession {
  id: string
  title: string
  lastMessage: string
  timestamp: Date
}

interface ChatSidebarProps {
  isOpen: boolean
  onClose: () => void
  chatSessions: ChatSession[]
  currentSessionId: string
  onNewChat: () => void
  onSelectSession: (sessionId: string) => void
  onDeleteSession: (sessionId: string) => void
}

export function ChatSidebar({
  isOpen,
  onClose,
  chatSessions,
  currentSessionId,
  onNewChat,
  onSelectSession,
  onDeleteSession
}: ChatSidebarProps) {
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editTitle, setEditTitle] = useState('')

  const handleEditStart = (session: ChatSession) => {
    setEditingId(session.id)
    setEditTitle(session.title)
  }

  const handleEditSave = () => {
    // TODO: Implement title update
    setEditingId(null)
    setEditTitle('')
  }

  const handleEditCancel = () => {
    setEditingId(null)
    setEditTitle('')
  }

  const formatTime = (date: Date) => {
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const minutes = Math.floor(diff / 60000)
    const hours = Math.floor(diff / 3600000)
    const days = Math.floor(diff / 86400000)

    if (minutes < 1) return 'Just now'
    if (minutes < 60) return `${minutes}m ago`
    if (hours < 24) return `${hours}h ago`
    if (days < 7) return `${days}d ago`
    return date.toLocaleDateString()
  }

  return (
    <>
      {/* Overlay for mobile */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black/50 z-40 lg:hidden" 
          onClick={onClose}
        />
      )}
      
      {/* Sidebar */}
      <div className={cn(
        "fixed lg:relative inset-y-0 left-0 z-50 w-80 bg-background border-r transform transition-transform duration-200 ease-in-out lg:translate-x-0",
        isOpen ? "translate-x-0" : "-translate-x-full"
      )}>
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b">
            <h2 className="text-lg font-semibold">Chats</h2>
            <div className="flex items-center space-x-2">
              <Button
                onClick={onNewChat}
                size="sm"
                variant="outline"
                className=""
              >
                <Plus className="h-4 w-4 mr-2" />
                New Chat
              </Button>
              <Button
                onClick={onClose}
                size="sm"
                variant="ghost"
                className="lg:hidden"
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Chat Sessions List */}
          <ScrollArea className="flex-1 p-2">
            <div className="space-y-2">
              {chatSessions.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <MessageSquare className="h-8 w-8 mx-auto mb-2 opacity-50" />
                  <p className="text-sm">No chats yet</p>
                  <p className="text-xs">Start a new conversation</p>
                </div>
              ) : (
                chatSessions.map((session) => (
                  <div
                    key={session.id}
                    className={cn(
                      "group relative p-3 rounded-lg border cursor-pointer transition-colors hover:bg-accent",
                      currentSessionId === session.id ? "bg-accent border-primary" : "border-border"
                    )}
                    onClick={() => onSelectSession(session.id)}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1 min-w-0">
                        {editingId === session.id ? (
                          <div className="flex items-center space-x-2">
                            <input
                              type="text"
                              value={editTitle}
                              onChange={(e) => setEditTitle(e.target.value)}
                              className="flex-1 px-2 py-1 text-sm border rounded"
                              onKeyDown={(e) => {
                                if (e.key === 'Enter') handleEditSave()
                                if (e.key === 'Escape') handleEditCancel()
                              }}
                              autoFocus
                            />
                          </div>
                        ) : (
                          <>
                            <h3 className="font-medium text-sm truncate">
                              {session.title}
                            </h3>
                            <p className="text-xs text-muted-foreground truncate mt-1">
                              {session.lastMessage || 'No messages yet'}
                            </p>
                            <div className="flex items-center justify-between mt-2">
                              <span className="text-xs text-muted-foreground">
                                {formatTime(session.timestamp)}
                              </span>
                              {currentSessionId === session.id && (
                                <Badge variant="default" className="text-xs">
                                  Active
                                </Badge>
                              )}
                            </div>
                          </>
                        )}
                      </div>
                      
                      {editingId !== session.id && (
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              variant="ghost"
                              size="sm"
                              className="opacity-0 group-hover:opacity-100 h-6 w-6 p-0"
                              onClick={(e) => e.stopPropagation()}
                            >
                              <MoreVertical className="h-3 w-3" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end" className="">
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                handleEditStart(session)
                              }}
                              className=""
                              inset={false}
                            >
                              <Edit className="mr-2 h-4 w-4" />
                              <span>Rename</span>
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={(e) => {
                                e.stopPropagation()
                                onDeleteSession(session.id)
                              }}
                              className="text-destructive"
                              inset={false}
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              <span>Delete</span>
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </ScrollArea>

          {/* Footer */}
          <div className="p-4 border-t">
            <div className="text-xs text-muted-foreground text-center">
              <p>GRYT AI Chat</p>
              <p>Powered by Advanced AI</p>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}