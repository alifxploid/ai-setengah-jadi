'use client'

import { useState, useEffect, useRef } from 'react'
import { ChatNavbar } from '@/components/chat/chat-navbar'
import { ChatSidebar } from '@/components/chat/chat-sidebar'
import { ChatArea } from '@/components/chat/chat-area'
import { ChatInput } from '@/components/chat/chat-input'
import { Button } from '@/components/ui/button'
import { Menu } from 'lucide-react'

interface Message {
  id: string
  content: string
  role: 'user' | 'assistant'
  timestamp: Date
}

interface ChatSession {
  id: string
  title: string
  lastMessage: string
  timestamp: Date
}

export default function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([])
  const [chatSessions, setChatSessions] = useState<ChatSession[]>([
    {
      id: '1',
      title: 'New Chat',
      lastMessage: 'Hello! How can I help you today?',
      timestamp: new Date()
    }
  ])
  const [currentSessionId, setCurrentSessionId] = useState('1')
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [isLoading, setIsLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSendMessage = async (content: string) => {
    if (!content.trim()) return

    const userMessage: Message = {
      id: Date.now().toString(),
      content: content.trim(),
      role: 'user',
      timestamp: new Date()
    }

    setMessages(prev => [...prev, userMessage])
    setIsLoading(true)

    try {
      // TODO: Implement API call to backend
      const response = await fetch('/api/chat/send', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          message: content,
          sessionId: currentSessionId
        }),
      })

      if (response.ok) {
        const data = await response.json()
        const assistantMessage: Message = {
          id: (Date.now() + 1).toString(),
          content: data.response || 'I\'m here to help! What would you like to know?',
          role: 'assistant',
          timestamp: new Date()
        }
        setMessages(prev => [...prev, assistantMessage])
      } else {
        throw new Error('Failed to send message')
      }
    } catch (error) {
      console.error('Error sending message:', error)
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: 'Sorry, I encountered an error. Please try again.',
        role: 'assistant',
        timestamp: new Date()
      }
      setMessages(prev => [...prev, errorMessage])
    } finally {
      setIsLoading(false)
    }
  }

  const handleNewChat = () => {
    const newSession: ChatSession = {
      id: Date.now().toString(),
      title: 'New Chat',
      lastMessage: '',
      timestamp: new Date()
    }
    setChatSessions(prev => [newSession, ...prev])
    setCurrentSessionId(newSession.id)
    setMessages([])
  }

  const handleSelectSession = (sessionId: string) => {
    setCurrentSessionId(sessionId)
    // TODO: Load messages for selected session
    setMessages([])
  }

  const handleDeleteSession = (sessionId: string) => {
    setChatSessions(prev => prev.filter(session => session.id !== sessionId))
    if (currentSessionId === sessionId && chatSessions.length > 1) {
      const remainingSessions = chatSessions.filter(session => session.id !== sessionId)
      setCurrentSessionId(remainingSessions[0].id)
    }
  }

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Navbar */}
      <ChatNavbar />
      
      <div className="flex flex-1 overflow-hidden">
        {/* Mobile sidebar toggle */}
        <div className="lg:hidden">
          {!sidebarOpen && (
            <Button
              variant="ghost"
              size="icon"
              className="fixed top-20 left-4 z-40"
              onClick={() => setSidebarOpen(true)}
            >
              <Menu className="h-5 w-5" />
            </Button>
          )}
        </div>

        {/* Sidebar */}
        <ChatSidebar
          isOpen={sidebarOpen}
          onClose={() => setSidebarOpen(false)}
          chatSessions={chatSessions}
          currentSessionId={currentSessionId}
          onNewChat={handleNewChat}
          onSelectSession={handleSelectSession}
          onDeleteSession={handleDeleteSession}
        />

        {/* Main Chat Area */}
        <div className="flex-1 flex flex-col">
          {/* Chat Messages */}
          <ChatArea
            messages={messages}
            isLoading={isLoading}
            messagesEndRef={messagesEndRef}
          />

          {/* Chat Input */}
          <ChatInput
            onSendMessage={handleSendMessage}
            disabled={isLoading}
          />
        </div>
      </div>
    </div>
  )
}