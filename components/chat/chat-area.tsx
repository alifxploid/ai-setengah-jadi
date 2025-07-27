'use client'

import { RefObject } from 'react'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Card } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Bot, User, Copy, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useState } from 'react'
import { cn } from '@/lib/utils'

interface Message {
  id: string
  content: string
  role: 'user' | 'assistant'
  timestamp: Date
}

interface ChatAreaProps {
  messages: Message[]
  isLoading: boolean
  messagesEndRef: RefObject<HTMLDivElement>
}

export function ChatArea({ messages, isLoading, messagesEndRef }: ChatAreaProps) {
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const handleCopy = async (content: string, messageId: string) => {
    try {
      await navigator.clipboard.writeText(content)
      setCopiedId(messageId)
      setTimeout(() => setCopiedId(null), 2000)
    } catch (error) {
      console.error('Failed to copy text:', error)
    }
  }

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    })
  }

  if (messages.length === 0 && !isLoading) {
    return (
      <div className="flex-1 flex items-center justify-center p-8">
        <div className="text-center max-w-md">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-primary/10 flex items-center justify-center">
            <Bot className="w-8 h-8 text-primary" />
          </div>
          <h3 className="text-lg font-semibold mb-2">Welcome to GRYT AI Chat</h3>
          <p className="text-muted-foreground mb-4">
            Start a conversation by typing a message below. I'm here to help you with any questions or tasks you have.
          </p>
          <div className="space-y-2 text-sm text-muted-foreground">
            <p>💡 Ask me anything</p>
            <p>🚀 Get instant responses</p>
            <p>🎯 Powered by advanced AI</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <ScrollArea className="flex-1 p-4">
      <div className="max-w-4xl mx-auto space-y-6">
        {messages.map((message) => (
          <div
            key={message.id}
            className={cn(
              "flex gap-4",
              message.role === 'user' ? "justify-end" : "justify-start"
            )}
          >
            {message.role === 'assistant' && (
              <Avatar className="w-8 h-8 mt-1">
                <AvatarImage src="" alt="AI" className="" />
                <AvatarFallback className="bg-primary text-primary-foreground">
                  <Bot className="w-4 h-4" />
                </AvatarFallback>
              </Avatar>
            )}
            
            <div className={cn(
              "flex flex-col max-w-[80%]",
              message.role === 'user' ? "items-end" : "items-start"
            )}>
              <Card className={cn(
                "p-4 relative group",
                message.role === 'user' 
                  ? "bg-primary text-primary-foreground" 
                  : "bg-muted"
              )}>
                <div className="whitespace-pre-wrap break-words">
                  {message.content}
                </div>
                
                {message.role === 'assistant' && (
                  <Button
                    variant="ghost"
                    size="sm"
                    className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity h-6 w-6 p-0"
                    onClick={() => handleCopy(message.content, message.id)}
                  >
                    {copiedId === message.id ? (
                      <Check className="w-3 h-3" />
                    ) : (
                      <Copy className="w-3 h-3" />
                    )}
                  </Button>
                )}
              </Card>
              
              <div className={cn(
                "flex items-center gap-2 mt-1 text-xs text-muted-foreground",
                message.role === 'user' ? "flex-row-reverse" : "flex-row"
              )}>
                <span>{message.role === 'user' ? 'You' : 'AI'}</span>
                <span>•</span>
                <span>{formatTime(message.timestamp)}</span>
              </div>
            </div>
            
            {message.role === 'user' && (
              <Avatar className="w-8 h-8 mt-1">
                <AvatarImage src="" alt="User" className="" />
                <AvatarFallback className="">
                  <User className="w-4 h-4" />
                </AvatarFallback>
              </Avatar>
            )}
          </div>
        ))}
        
        {isLoading && (
          <div className="flex gap-4 justify-start">
            <Avatar className="w-8 h-8 mt-1">
              <AvatarImage src="" alt="AI" className="" />
              <AvatarFallback className="bg-primary text-primary-foreground">
                <Bot className="w-4 h-4" />
              </AvatarFallback>
            </Avatar>
            
            <div className="flex flex-col max-w-[80%] items-start">
              <Card className="p-4 bg-muted">
                <div className="space-y-2">
                  <Skeleton className="h-4 w-[200px]" />
                  <Skeleton className="h-4 w-[160px]" />
                  <Skeleton className="h-4 w-[180px]" />
                </div>
              </Card>
              
              <div className="flex items-center gap-2 mt-1 text-xs text-muted-foreground">
                <span>AI</span>
                <span>•</span>
                <span>Typing...</span>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>
    </ScrollArea>
  )
}