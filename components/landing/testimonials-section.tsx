'use client'

import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Star, Quote } from 'lucide-react'

const testimonials = [
  {
    name: 'Sarah Chen',
    role: 'Senior Developer',
    company: 'TechCorp',
    avatar: '/avatars/sarah.jpg',
    content: 'GRYT AI has completely transformed our development workflow. We\'re shipping features 3x faster than before.',
    rating: 5,
    badge: 'Verified User'
  },
  {
    name: 'Marcus Rodriguez',
    role: 'CTO',
    company: 'StartupXYZ',
    avatar: '/avatars/marcus.jpg',
    content: 'The AI code generation is incredibly accurate. It feels like having a senior developer pair programming with you 24/7.',
    rating: 5,
    badge: 'Enterprise'
  },
  {
    name: 'Emily Johnson',
    role: 'Full Stack Developer',
    company: 'InnovateLab',
    avatar: '/avatars/emily.jpg',
    content: 'I was skeptical at first, but GRYT AI has saved me countless hours. The deployment speed is unmatched.',
    rating: 5,
    badge: 'Pro User'
  },
  {
    name: 'David Kim',
    role: 'Lead Engineer',
    company: 'CloudScale',
    avatar: '/avatars/david.jpg',
    content: 'The security features and compliance tools give us peace of mind when handling sensitive data.',
    rating: 5,
    badge: 'Enterprise'
  },
  {
    name: 'Lisa Wang',
    role: 'Product Manager',
    company: 'DigitalFlow',
    avatar: '/avatars/lisa.jpg',
    content: 'Our team productivity has increased dramatically. The analytics insights help us make better decisions.',
    rating: 5,
    badge: 'Verified User'
  },
  {
    name: 'Alex Thompson',
    role: 'Freelance Developer',
    company: 'Independent',
    avatar: '/avatars/alex.jpg',
    content: 'As a freelancer, GRYT AI helps me deliver high-quality projects faster, which means happier clients.',
    rating: 5,
    badge: 'Pro User'
  }
]

function StarRating({ rating }: { rating: number }) {
  return (
    <div className="flex items-center space-x-1">
      {[...Array(5)].map((_, i) => (
        <Star
          key={i}
          className={`h-4 w-4 ${
            i < rating ? 'fill-yellow-400 text-yellow-400' : 'text-muted-foreground'
          }`}
        />
      ))}
    </div>
  )
}

export function TestimonialsSection() {
  return (
    <section className="py-20 sm:py-32">
      <div className="container">
        <div className="mx-auto max-w-2xl text-center mb-16">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl mb-4">
            Loved by developers worldwide
          </h2>
          <p className="text-lg text-muted-foreground">
            Join thousands of developers who are building amazing applications with GRYT AI.
          </p>
        </div>
        
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {testimonials.map((testimonial, index) => (
            <Card key={index} className="relative overflow-hidden border-0 bg-card/50 backdrop-blur-sm transition-all duration-300 hover:bg-card hover:shadow-lg">
              <CardContent className="p-6">
                <div className="flex items-start justify-between mb-4">
                  <Quote className="h-8 w-8 text-primary/20" />
                  <Badge variant="secondary" className="text-xs">
                    {testimonial.badge}
                  </Badge>
                </div>
                
                <blockquote className="text-sm leading-relaxed mb-6">
                  "{testimonial.content}"
                </blockquote>
                
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <Avatar className="h-10 w-10">
                      <AvatarImage src={testimonial.avatar} alt={testimonial.name} className="object-cover" />
                      <AvatarFallback className="bg-primary/10 text-primary font-semibold">
                        {testimonial.name.split(' ').map(n => n[0]).join('')}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <p className="font-semibold text-sm">{testimonial.name}</p>
                      <p className="text-xs text-muted-foreground">
                        {testimonial.role} at {testimonial.company}
                      </p>
                    </div>
                  </div>
                  <StarRating rating={testimonial.rating} />
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
        
        <div className="mt-16 text-center">
          <div className="inline-flex items-center space-x-2 text-sm text-muted-foreground">
            <div className="flex items-center space-x-1">
              {[...Array(5)].map((_, i) => (
                <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
              ))}
            </div>
            <span>4.9/5 from 2,000+ reviews</span>
          </div>
        </div>
      </div>
    </section>
  )
}