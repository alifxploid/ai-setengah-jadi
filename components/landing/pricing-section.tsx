'use client'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Check, Star } from 'lucide-react'

const pricingPlans = [
  {
    name: 'Starter',
    description: 'Perfect for individuals and small projects',
    price: 'Free',
    period: 'forever',
    features: [
      '5 AI-generated projects per month',
      'Basic code completion',
      'Community support',
      '1GB storage',
      'Standard deployment speed'
    ],
    cta: 'Get Started',
    popular: false
  },
  {
    name: 'Pro',
    description: 'Ideal for growing teams and businesses',
    price: '$29',
    period: 'per month',
    features: [
      'Unlimited AI-generated projects',
      'Advanced code completion',
      'Priority support',
      '100GB storage',
      'Lightning-fast deployment',
      'Custom integrations',
      'Advanced analytics'
    ],
    cta: 'Start Free Trial',
    popular: true
  },
  {
    name: 'Enterprise',
    description: 'For large organizations with custom needs',
    price: 'Custom',
    period: 'contact us',
    features: [
      'Everything in Pro',
      'Dedicated account manager',
      'Custom AI model training',
      'Unlimited storage',
      'SLA guarantee',
      'On-premise deployment',
      'Advanced security features',
      'Custom integrations'
    ],
    cta: 'Contact Sales',
    popular: false
  }
]

export function PricingSection() {
  return (
    <section id="pricing" className="py-20 sm:py-32 bg-muted/30">
      <div className="container">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl mb-4">
            Simple, transparent pricing
          </h2>
          <p className="text-lg text-muted-foreground">
            Choose the perfect plan for your needs. Upgrade or downgrade at any time.
          </p>
        </div>
        
        <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
          {pricingPlans.map((plan, index) => (
            <Card 
              key={index} 
              className={`relative overflow-hidden transition-all duration-300 hover:shadow-lg ${
                plan.popular 
                  ? 'border-primary shadow-lg scale-105 bg-card' 
                  : 'border-border bg-card/50 backdrop-blur-sm'
              }`}
            >
              {plan.popular && (
                <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2">
                  <Badge variant="default" className="bg-primary text-primary-foreground px-4 py-1">
                    <Star className="w-3 h-3 mr-1" />
                    Most Popular
                  </Badge>
                </div>
              )}
              
              <CardHeader className="text-center pb-8 pt-8">
                <CardTitle className="text-2xl font-bold">{plan.name}</CardTitle>
                <CardDescription className="text-base mt-2">
                  {plan.description}
                </CardDescription>
                <div className="mt-6">
                  <div className="flex items-baseline justify-center">
                    <span className="text-4xl font-bold tracking-tight">
                      {plan.price}
                    </span>
                    {plan.price !== 'Free' && plan.price !== 'Custom' && (
                      <span className="text-muted-foreground ml-1">/{plan.period.split(' ')[0]}</span>
                    )}
                  </div>
                  <p className="text-sm text-muted-foreground mt-1">
                    {plan.period}
                  </p>
                </div>
              </CardHeader>
              
              <CardContent className="pt-0">
                <ul className="space-y-3 mb-8">
                  {plan.features.map((feature, featureIndex) => (
                    <li key={featureIndex} className="flex items-start">
                      <Check className="h-4 w-4 text-primary mt-0.5 mr-3 flex-shrink-0" />
                      <span className="text-sm">{feature}</span>
                    </li>
                  ))}
                </ul>
                
                <Button 
                  className="w-full" 
                  variant={plan.popular ? 'default' : 'outline'}
                  size="lg"
                >
                  {plan.cta}
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
        
        <div className="mt-16 text-center">
          <p className="text-sm text-muted-foreground">
            All plans include a 14-day free trial. No credit card required.
          </p>
        </div>
      </div>
    </section>
  )
}