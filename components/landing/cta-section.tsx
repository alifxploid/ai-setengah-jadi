'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowRight, Sparkles, Users, Zap } from 'lucide-react'

export function CTASection() {
  return (
    <section className="py-20 sm:py-32 bg-gradient-to-br from-primary/5 via-purple-500/5 to-pink-500/5">
      <div className="container">
        <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-primary via-purple-600 to-pink-600 px-8 py-16 sm:px-16 sm:py-24">
          {/* Background decoration */}
          <div className="absolute inset-0 bg-gradient-to-br from-white/10 to-transparent" />
          <div className="absolute -top-40 -right-40 h-80 w-80 rounded-full bg-white/10 blur-3xl" />
          <div className="absolute -bottom-40 -left-40 h-80 w-80 rounded-full bg-white/10 blur-3xl" />
          
          <div className="relative mx-auto max-w-4xl text-center">
            {/* Badge */}
            <div className="mb-8 flex justify-center">
              <Badge variant="secondary" className="bg-white/20 text-white border-white/30 px-4 py-2 text-sm font-medium backdrop-blur-sm">
                <Sparkles className="mr-2 h-4 w-4" />
                Join 50,000+ Developers
              </Badge>
            </div>
            
            {/* Main heading */}
            <h2 className="mb-6 text-4xl font-bold tracking-tight text-white sm:text-5xl lg:text-6xl">
              Ready to transform your development workflow?
            </h2>
            
            {/* Subheading */}
            <p className="mb-8 text-lg text-white/90 sm:text-xl lg:text-2xl">
              Start building smarter today with GRYT AI. No credit card required.
            </p>
            
            {/* CTA Buttons */}
            <div className="mb-12 flex flex-col gap-4 sm:flex-row sm:justify-center">
              <Button 
                size="lg" 
                variant="secondary" 
                className="h-12 px-8 text-base bg-white text-primary hover:bg-white/90" 
                asChild
              >
                <Link href="/auth/register">
                  Start Free Trial
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
              <Button 
                size="lg" 
                variant="outline" 
                className="h-12 px-8 text-base border-white/30 text-white hover:bg-white/10" 
                asChild
              >
                <Link href="#demo">
                  Schedule Demo
                </Link>
              </Button>
            </div>
            
            {/* Stats */}
            <div className="grid grid-cols-1 gap-8 sm:grid-cols-3">
              <div className="flex flex-col items-center space-y-2">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white/20 backdrop-blur-sm">
                  <Users className="h-6 w-6 text-white" />
                </div>
                <div className="text-2xl font-bold text-white">50K+</div>
                <div className="text-sm text-white/80">Active Developers</div>
              </div>
              
              <div className="flex flex-col items-center space-y-2">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white/20 backdrop-blur-sm">
                  <Zap className="h-6 w-6 text-white" />
                </div>
                <div className="text-2xl font-bold text-white">1M+</div>
                <div className="text-sm text-white/80">Projects Deployed</div>
              </div>
              
              <div className="flex flex-col items-center space-y-2">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white/20 backdrop-blur-sm">
                  <Sparkles className="h-6 w-6 text-white" />
                </div>
                <div className="text-2xl font-bold text-white">99.9%</div>
                <div className="text-sm text-white/80">Uptime Guarantee</div>
              </div>
            </div>
          </div>
        </div>
        
        {/* Trust indicators */}
        <div className="mt-16 text-center">
          <p className="text-sm text-muted-foreground mb-8">
            Trusted by teams at leading companies worldwide
          </p>
          <div className="flex flex-wrap items-center justify-center gap-8 opacity-60">
            <div className="text-2xl font-bold text-muted-foreground">Google</div>
            <div className="text-2xl font-bold text-muted-foreground">Microsoft</div>
            <div className="text-2xl font-bold text-muted-foreground">Amazon</div>
            <div className="text-2xl font-bold text-muted-foreground">Netflix</div>
            <div className="text-2xl font-bold text-muted-foreground">Spotify</div>
          </div>
        </div>
      </div>
    </section>
  )
}