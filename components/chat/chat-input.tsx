'use client'

import { useState, useRef, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card } from '@/components/ui/card'
import { Send, Paperclip, Mic, Square } from 'lucide-react'
import { cn } from '@/lib/utils'

interface ChatInputProps {
  onSendMessage: (message: string) => void
  disabled?: boolean
}

export function ChatInput({ onSendMessage, disabled = false }: ChatInputProps) {
  const [message, setMessage] = useState('')
  const [isRecording, setIsRecording] = useState(false)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (message.trim() && !disabled) {
      onSendMessage(message)
      setMessage('')
      resetTextareaHeight()
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSubmit(e)
    }
  }

  const resetTextareaHeight = () => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
    }
  }

  const adjustTextareaHeight = () => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      const scrollHeight = textareaRef.current.scrollHeight
      const maxHeight = 200 // Maximum height in pixels
      textareaRef.current.style.height = `${Math.min(scrollHeight, maxHeight)}px`
    }
  }

  useEffect(() => {
    adjustTextareaHeight()
  }, [message])

  const handleVoiceRecord = () => {
    if (isRecording) {
      // Stop recording
      setIsRecording(false)
      // TODO: Implement voice recording stop
    } else {
      // Start recording
      setIsRecording(true)
      // TODO: Implement voice recording start
    }
  }

  const handleFileAttach = () => {
    // TODO: Implement file attachment
    console.log('File attachment clicked')
  }

  return (
    <div className="border-t bg-background p-4">
      <div className="max-w-4xl mx-auto">
        <Card className="p-4">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="flex items-end gap-2">
              {/* File Attachment Button */}
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="mb-2"
                onClick={handleFileAttach}
                disabled={disabled}
              >
                <Paperclip className="h-4 w-4" />
              </Button>

              {/* Message Input */}
              <div className="flex-1 relative">
                <Textarea
                  ref={textareaRef}
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder={disabled ? "AI is thinking..." : "Type your message here... (Press Enter to send, Shift+Enter for new line)"}
                  disabled={disabled}
                  className={cn(
                    "min-h-[60px] max-h-[200px] resize-none pr-12",
                    "focus:ring-2 focus:ring-primary focus:border-transparent"
                  )}
                  rows={1}
                />
                
                {/* Character count */}
                <div className="absolute bottom-2 right-2 text-xs text-muted-foreground">
                  {message.length}/2000
                </div>
              </div>

              {/* Voice Recording Button */}
              <Button
                type="button"
                variant={isRecording ? "destructive" : "ghost"}
                size="sm"
                className="mb-2"
                onClick={handleVoiceRecord}
                disabled={disabled}
              >
                {isRecording ? (
                  <Square className="h-4 w-4" />
                ) : (
                  <Mic className="h-4 w-4" />
                )}
              </Button>

              {/* Send Button */}
              <Button
                type="submit"
                disabled={!message.trim() || disabled}
                size="sm"
                variant="default"
                className="mb-2"
              >
                <Send className="h-4 w-4" />
              </Button>
            </div>

            {/* Quick Actions */}
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setMessage("Explain this concept to me: ")}
                disabled={disabled}
                className="text-xs"
              >
                💡 Explain
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setMessage("Help me write code for: ")}
                disabled={disabled}
                className="text-xs"
              >
                💻 Code Help
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setMessage("Summarize this for me: ")}
                disabled={disabled}
                className="text-xs"
              >
                📝 Summarize
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setMessage("What are the pros and cons of: ")}
                disabled={disabled}
                className="text-xs"
              >
                ⚖️ Analyze
              </Button>
            </div>
          </form>
        </Card>
        
        {/* Footer Info */}
        <div className="text-center mt-2">
          <p className="text-xs text-muted-foreground">
            GRYT AI can make mistakes. Please verify important information.
          </p>
        </div>
      </div>
    </div>
  )
}