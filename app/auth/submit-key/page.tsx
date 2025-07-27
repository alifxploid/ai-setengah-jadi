'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Key, Zap, CheckCircle, AlertCircle, LogOut, User, Shield, Clock, Mail } from 'lucide-react'

export default function SubmitKeyPage() {
  const [accessKey, setAccessKey] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')
    setSuccess('')

    if (!accessKey.trim()) {
      setError('Access key is required')
      setIsLoading(false)
      return
    }

    // Validate access key format (example: should be 32 characters alphanumeric)
    if (accessKey.trim().length < 16) {
      setError('Invalid access key format')
      setIsLoading(false)
      return
    }

    try {
      // TODO: Implement access key validation API call
      const response = await fetch('/api/auth/validate-access-key', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          accessKey: accessKey.trim()
        }),
      })

      if (response.ok) {
        const data = await response.json()
        setSuccess('Access key validated successfully! Redirecting to chat...')
        setAccessKey('')
        
        // Store access token in localStorage or cookie
        localStorage.setItem('gryt_access_token', data.token)
        
        // Redirect to chat page after 2 seconds
        setTimeout(() => {
          window.location.href = '/chat'
        }, 2000)
      } else {
        const data = await response.json()
        setError(data.message || 'Invalid access key')
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
            <h1 className="text-3xl font-bold mb-2">Access Key Required</h1>
            <p className="text-muted-foreground">
              Enter the access key provided by admin to access GRYT AI Chat
            </p>
          </div>

          <div className="max-w-md mx-auto">
            {/* Submit Access Key */}
            <Card className="">
              <CardHeader className="text-center">
                <div className="h-16 w-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Shield className="h-8 w-8 text-primary" />
                </div>
                <CardTitle className="text-xl">
                  Enter Access Key
                </CardTitle>
                <CardDescription className="">
                  Use the access key provided by your administrator
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
                    <Label htmlFor="accessKey" className="">Access Key</Label>
                    <Input
                      id="accessKey"
                      type="password"
                      placeholder="Enter your access key..."
                      value={accessKey}
                      onChange={(e) => setAccessKey(e.target.value)}
                      className="font-mono text-center tracking-wider"
                      required
                    />
                    <p className="text-xs text-muted-foreground text-center">
                      Contact your administrator if you don't have an access key
                    </p>
                  </div>

                  <Button type="submit" className="w-full" variant="default" size="default" disabled={isLoading}>
                    {isLoading ? 'Validating...' : 'Access Chat'}
                  </Button>
                </form>
              </CardContent>
            </Card>
          </div>



          {/* Info Section */}
          <Card className="bg-muted/50 max-w-2xl mx-auto">
            <CardContent className="pt-6">
              <div className="grid gap-6 md:grid-cols-3">
                <div className="text-center">
                  <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-3">
                    <Shield className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-medium mb-1">Secure Access</h3>
                  <p className="text-sm text-muted-foreground">
                    Access keys are encrypted and validated securely
                  </p>
                </div>
                <div className="text-center">
                  <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-3">
                    <Clock className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-medium mb-1">Instant Access</h3>
                  <p className="text-sm text-muted-foreground">
                    Valid keys provide immediate chat access
                  </p>
                </div>
                <div className="text-center">
                  <div className="h-12 w-12 bg-primary/10 rounded-lg flex items-center justify-center mx-auto mb-3">
                    <Mail className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-medium mb-1">Need Help?</h3>
                  <p className="text-sm text-muted-foreground">
                    Contact admin for access key assistance
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