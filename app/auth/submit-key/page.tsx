'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Textarea } from '@/components/ui/textarea'
import { Key, Zap, CheckCircle, AlertCircle, LogOut, User } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

export default function SubmitKeyPage() {
  const [apiKey, setApiKey] = useState('')
  const [description, setDescription] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [submittedKeys, setSubmittedKeys] = useState([
    {
      id: '1',
      name: 'OpenAI GPT-4',
      status: 'active',
      submittedAt: '2024-01-15'
    },
    {
      id: '2', 
      name: 'Anthropic Claude',
      status: 'pending',
      submittedAt: '2024-01-14'
    }
  ])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')
    setSuccess('')

    if (!apiKey.trim()) {
      setError('API Key is required')
      setIsLoading(false)
      return
    }

    try {
      // TODO: Implement submit key API call
      const response = await fetch('/api/keys/submit', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          apiKey: apiKey.trim(),
          description: description.trim()
        }),
      })

      if (response.ok) {
        const data = await response.json()
        setSuccess('API Key submitted successfully! Redirecting to chat...')
        setApiKey('')
        setDescription('')
        // Add to submitted keys list
        setSubmittedKeys(prev => [{
          id: Date.now().toString(),
          name: description || 'API Key',
          status: 'pending',
          submittedAt: new Date().toISOString().split('T')[0]
        }, ...prev])
        
        // Redirect to chat page after 2 seconds
        setTimeout(() => {
          window.location.href = '/chat'
        }, 2000)
      } else {
        const data = await response.json()
        setError(data.message || 'Failed to submit API key')
      }
    } catch (err) {
      setError('Something went wrong. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleLogout = () => {
    // TODO: Implement logout
    window.location.href = '/auth/login'
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'active':
        return <Badge variant="default" className="bg-green-500 hover:bg-green-600">Active</Badge>
      case 'pending':
        return <Badge variant="secondary" className="">Pending</Badge>
      case 'rejected':
        return <Badge variant="destructive" className="">Rejected</Badge>
      default:
        return <Badge variant="outline" className="">Unknown</Badge>
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
      {/* Header */}
      <header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link href="/" className="flex items-center space-x-2">
              <Zap className="h-6 w-6 text-primary" />
              <span className="font-bold text-xl">GRYT AI</span>
            </Link>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                <User className="h-4 w-4" />
                <span>Welcome back!</span>
              </div>
              <Button variant="outline" size="sm" onClick={handleLogout} className="">
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto space-y-8">
          {/* Page Title */}
          <div className="text-center">
            <h1 className="text-3xl font-bold mb-2">API Key Management</h1>
            <p className="text-muted-foreground">
              Submit your API keys to unlock the full potential of GRYT AI
            </p>
          </div>

          <div className="grid gap-8 md:grid-cols-2">
            {/* Submit New Key */}
            <Card className="">
              <CardHeader className="">
                <CardTitle className="flex items-center space-x-2">
                  <Key className="h-5 w-5" />
                  <span>Submit New API Key</span>
                </CardTitle>
                <CardDescription className="">
                  Add a new API key to expand your AI capabilities
                </CardDescription>
              </CardHeader>
              <CardContent className="">
                <form onSubmit={handleSubmit} className="space-y-4">
                  {error && (
                    <Alert variant="destructive" className="">
                      <AlertCircle className="h-4 w-4" />
                      <AlertDescription className="">{error}</AlertDescription>
                    </Alert>
                  )}
                  
                  {success && (
                    <Alert variant="default" className="border-green-200 bg-green-50 text-green-800">
                      <CheckCircle className="h-4 w-4" />
                      <AlertDescription className="">{success}</AlertDescription>
                    </Alert>
                  )}

                  <div className="space-y-2">
                    <Label htmlFor="description" className="">Description (Optional)</Label>
                    <Input
                      id="description"
                      type="text"
                      placeholder="e.g., OpenAI GPT-4, Anthropic Claude"
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                      className=""
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="apiKey" className="">API Key</Label>
                    <Textarea
                      id="apiKey"
                      placeholder="Paste your API key here..."
                      value={apiKey}
                      onChange={(e) => setApiKey(e.target.value)}
                      className="min-h-[100px] font-mono text-sm"
                      required
                    />
                    <p className="text-xs text-muted-foreground">
                      Your API key will be encrypted and stored securely
                    </p>
                  </div>

                  <Button type="submit" className="w-full" variant="default" size="default" disabled={isLoading}>
                    {isLoading ? 'Submitting...' : 'Submit API Key'}
                  </Button>
                </form>
              </CardContent>
            </Card>

            {/* Submitted Keys */}
            <Card className="">
              <CardHeader className="">
                <CardTitle className="">Your API Keys</CardTitle>
                <CardDescription className="">
                  Manage your submitted API keys
                </CardDescription>
              </CardHeader>
              <CardContent className="">
                <div className="space-y-4">
                  {submittedKeys.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                      <Key className="h-12 w-12 mx-auto mb-4 opacity-50" />
                      <p>No API keys submitted yet</p>
                    </div>
                  ) : (
                    submittedKeys.map((key) => (
                      <div key={key.id} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex-1">
                          <h4 className="font-medium">{key.name}</h4>
                          <p className="text-sm text-muted-foreground">
                            Submitted on {key.submittedAt}
                          </p>
                        </div>
                        <div className="flex items-center space-x-2">
                          {getStatusBadge(key.status)}
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Info Section */}
          <Card className="bg-muted/50">
            <CardContent className="pt-6">
              <div className="grid gap-4 md:grid-cols-3">
                <div className="text-center">
                  <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-2">
                    <Key className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-medium mb-1">Secure Storage</h3>
                  <p className="text-sm text-muted-foreground">
                    All API keys are encrypted and stored securely
                  </p>
                </div>
                <div className="text-center">
                  <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-2">
                    <CheckCircle className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-medium mb-1">Quick Approval</h3>
                  <p className="text-sm text-muted-foreground">
                    Most keys are approved within 24 hours
                  </p>
                </div>
                <div className="text-center">
                  <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-2">
                    <Zap className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-medium mb-1">Enhanced Features</h3>
                  <p className="text-sm text-muted-foreground">
                    Unlock advanced AI capabilities
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}