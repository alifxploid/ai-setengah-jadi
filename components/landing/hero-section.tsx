'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowRight, Sparkles, Zap, Code, Brain } from 'lucide-react'

export function HeroSection() {
  return (
    <section className="relative overflow-hidden bg-gradient-to-br from-background via-background to-muted/20 py-20 sm:py-32">
      {/* Background decoration */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
          <div className="h-[800px] w-[800px] rounded-full bg-gradient-to-r from-primary/10 via-purple-500/10 to-pink-500/10 blur-3xl" />
        </div>
      </div>
      
      <div className="container relative">
        <div className="mx-auto max-w-4xl text-center">
          {/* Badge */}
          <div className="mb-8 flex justify-center">
            <Badge variant="secondary" className="px-4 py-2 text-sm font-medium">
              <Sparkles className="mr-2 h-4 w-4" />
              Powered by Advanced AI Technology
            </Badge>
          </div>
          
          {/* Main heading */}
          <h1 className="mb-6 text-4xl font-bold tracking-tight text-foreground sm:text-6xl lg:text-7xl">
            Build Smarter with{' '}
            <span className="bg-gradient-to-r from-primary via-purple-500 to-pink-500 bg-clip-text text-transparent">
              GRYT AI
            </span>
          </h1>
          
          {/* Subheading */}
          <p className="mb-8 text-lg text-muted-foreground sm:text-xl lg:text-2xl">
            Revolutionary AI-powered development platform that transforms your ideas into reality.
            Code faster, deploy smarter, and scale effortlessly with our intelligent tools.
          </p>
          
          {/* CTA Buttons */}
          <div className="mb-12 flex flex-col gap-4 sm:flex-row sm:justify-center">
            <Button size="lg" variant="default" className="h-12 px-8 text-base" asChild>
              <Link href="/auth/register">
                Get Started Free
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button size="lg" variant="outline" className="h-12 px-8 text-base" asChild>
              <Link href="#demo">
                Watch Demo
              </Link>
            </Button>
          </div>
          
          {/* Feature highlights */}
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-3">
            <div className="flex flex-col items-center space-y-2">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <Code className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold">Smart Code Generation</h3>
              <p className="text-sm text-muted-foreground">AI-powered code completion and generation</p>
            </div>
            
            <div className="flex flex-col items-center space-y-2">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <Zap className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold">Lightning Fast</h3>
              <p className="text-sm text-muted-foreground">Deploy and scale in seconds, not hours</p>
            </div>
            
            <div className="flex flex-col items-center space-y-2">
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                <Brain className="h-6 w-6 text-primary" />
              </div>
              <h3 className="font-semibold">Intelligent Insights</h3>
              <p className="text-sm text-muted-foreground">Advanced analytics and optimization</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}