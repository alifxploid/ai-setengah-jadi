'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { 
  Code2, 
  Zap, 
  Shield, 
  Globe, 
  Database, 
  Cpu, 
  GitBranch, 
  Smartphone,
  Cloud,
  BarChart3,
  Lock,
  Rocket
} from 'lucide-react'

const features = [
  {
    icon: Code2,
    title: 'AI Code Generation',
    description: 'Generate high-quality code instantly with our advanced AI models trained on millions of repositories.',
    badge: 'Core Feature'
  },
  {
    icon: Zap,
    title: 'Lightning Deployment',
    description: 'Deploy your applications in seconds with our optimized CI/CD pipeline and edge computing.',
    badge: 'Performance'
  },
  {
    icon: Shield,
    title: 'Enterprise Security',
    description: 'Bank-grade security with end-to-end encryption, SOC 2 compliance, and advanced threat protection.',
    badge: 'Security'
  },
  {
    icon: Globe,
    title: 'Global CDN',
    description: 'Serve your applications worldwide with our distributed network of edge servers.',
    badge: 'Infrastructure'
  },
  {
    icon: Database,
    title: 'Smart Database',
    description: 'Intelligent database optimization with automatic scaling and performance tuning.',
    badge: 'Database'
  },
  {
    icon: Cpu,
    title: 'Auto Scaling',
    description: 'Automatically scale your resources based on demand with intelligent load balancing.',
    badge: 'Scaling'
  },
  {
    icon: GitBranch,
    title: 'Version Control',
    description: 'Advanced Git integration with automated testing, code review, and deployment workflows.',
    badge: 'DevOps'
  },
  {
    icon: Smartphone,
    title: 'Mobile First',
    description: 'Build responsive applications that work perfectly across all devices and platforms.',
    badge: 'Mobile'
  },
  {
    icon: Cloud,
    title: 'Cloud Native',
    description: 'Built for the cloud with microservices architecture and containerized deployments.',
    badge: 'Cloud'
  },
  {
    icon: BarChart3,
    title: 'Real-time Analytics',
    description: 'Monitor your applications with comprehensive analytics and performance insights.',
    badge: 'Analytics'
  },
  {
    icon: Lock,
    title: 'Data Privacy',
    description: 'Your data stays private with zero-knowledge architecture and GDPR compliance.',
    badge: 'Privacy'
  },
  {
    icon: Rocket,
    title: 'Innovation Lab',
    description: 'Access cutting-edge features and experimental tools before they go mainstream.',
    badge: 'Innovation'
  }
]

export function FeaturesSection() {
  return (
    <section id="features" className="py-20 sm:py-32">
      <div className="container">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl mb-4">
            Everything you need to build amazing apps
          </h2>
          <p className="text-lg text-muted-foreground">
            Powerful features designed to accelerate your development workflow and deliver exceptional user experiences.
          </p>
        </div>
        
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {features.map((feature, index) => {
            const Icon = feature.icon
            return (
              <Card key={index} className="group relative overflow-hidden border-0 bg-card/50 backdrop-blur-sm transition-all duration-300 hover:bg-card hover:shadow-lg">
                <CardHeader className="pb-4">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 group-hover:bg-primary/20 transition-colors">
                      <Icon className="h-5 w-5 text-primary" />
                    </div>
                    <Badge variant="secondary" className="text-xs">
                      {feature.badge}
                    </Badge>
                  </div>
                  <CardTitle className="text-lg">{feature.title}</CardTitle>
                </CardHeader>
                <CardContent className="pt-0">
                  <CardDescription className="text-sm leading-relaxed">
                    {feature.description}
                  </CardDescription>
                </CardContent>
              </Card>
            )
          })}
        </div>
      </div>
    </section>
  )
}