'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { 
  Zap, 
  Github, 
  Twitter, 
  Linkedin, 
  Youtube,
  Mail,
  MapPin,
  Phone
} from 'lucide-react'

const footerLinks = {
  product: [
    { name: 'Features', href: '#features' },
    { name: 'Pricing', href: '#pricing' },
    { name: 'API Documentation', href: '/docs' },
    { name: 'Integrations', href: '/integrations' },
    { name: 'Changelog', href: '/changelog' }
  ],
  company: [
    { name: 'About Us', href: '/about' },
    { name: 'Careers', href: '/careers' },
    { name: 'Blog', href: '/blog' },
    { name: 'Press Kit', href: '/press' },
    { name: 'Contact', href: '/contact' }
  ],
  resources: [
    { name: 'Help Center', href: '/help' },
    { name: 'Community', href: '/community' },
    { name: 'Tutorials', href: '/tutorials' },
    { name: 'Status Page', href: '/status' },
    { name: 'System Requirements', href: '/requirements' }
  ],
  legal: [
    { name: 'Privacy Policy', href: '/privacy' },
    { name: 'Terms of Service', href: '/terms' },
    { name: 'Cookie Policy', href: '/cookies' },
    { name: 'GDPR', href: '/gdpr' },
    { name: 'Security', href: '/security' }
  ]
}

const socialLinks = [
  { name: 'GitHub', href: 'https://github.com/gryt-ai', icon: Github },
  { name: 'Twitter', href: 'https://twitter.com/gryt_ai', icon: Twitter },
  { name: 'LinkedIn', href: 'https://linkedin.com/company/gryt-ai', icon: Linkedin },
  { name: 'YouTube', href: 'https://youtube.com/@gryt-ai', icon: Youtube }
]

export function Footer() {
  return (
    <footer className="bg-muted/30 border-t">
      <div className="container">
        {/* Main footer content */}
        <div className="py-16">
          <div className="grid grid-cols-1 gap-8 lg:grid-cols-6">
            {/* Brand section */}
            <div className="lg:col-span-2">
              <Link href="/" className="flex items-center space-x-2 mb-4">
                <Zap className="h-8 w-8 text-primary" />
                <span className="font-bold text-2xl">GRYT AI</span>
              </Link>
              <p className="text-muted-foreground mb-6 max-w-sm">
                Revolutionary AI-powered development platform that transforms your ideas into reality. 
                Build faster, deploy smarter, scale effortlessly.
              </p>
              
              {/* Newsletter signup */}
              <div className="space-y-4">
                <h4 className="font-semibold">Stay updated</h4>
                <div className="flex space-x-2">
                  <Input 
                    placeholder="Enter your email" 
                    type="email" 
                    className="max-w-sm"
                  />
                  <Button variant="default" size="default" className="px-4">
                    Subscribe
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  Get the latest updates and news delivered to your inbox.
                </p>
              </div>
            </div>
            
            {/* Links sections */}
            <div className="lg:col-span-4">
              <div className="grid grid-cols-2 gap-8 sm:grid-cols-4">
                <div>
                  <h4 className="font-semibold mb-4">Product</h4>
                  <ul className="space-y-3">
                    {footerLinks.product.map((link) => (
                      <li key={link.name}>
                        <Link 
                          href={link.href} 
                          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {link.name}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </div>
                
                <div>
                  <h4 className="font-semibold mb-4">Company</h4>
                  <ul className="space-y-3">
                    {footerLinks.company.map((link) => (
                      <li key={link.name}>
                        <Link 
                          href={link.href} 
                          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {link.name}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </div>
                
                <div>
                  <h4 className="font-semibold mb-4">Resources</h4>
                  <ul className="space-y-3">
                    {footerLinks.resources.map((link) => (
                      <li key={link.name}>
                        <Link 
                          href={link.href} 
                          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {link.name}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </div>
                
                <div>
                  <h4 className="font-semibold mb-4">Legal</h4>
                  <ul className="space-y-3">
                    {footerLinks.legal.map((link) => (
                      <li key={link.name}>
                        <Link 
                          href={link.href} 
                          className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                        >
                          {link.name}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <Separator className="my-8" />
        
        {/* Bottom footer */}
        <div className="py-8">
          <div className="flex flex-col items-center justify-between space-y-4 sm:flex-row sm:space-y-0">
            <div className="flex flex-col items-center space-y-2 sm:flex-row sm:space-y-0 sm:space-x-4">
              <p className="text-sm text-muted-foreground">
                © 2024 GRYT AI. All rights reserved.
              </p>
              <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                <div className="flex items-center space-x-1">
                  <MapPin className="h-3 w-3" />
                  <span>San Francisco, CA</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Mail className="h-3 w-3" />
                  <span>hello@gryt.ai</span>
                </div>
              </div>
            </div>
            
            {/* Social links */}
            <div className="flex items-center space-x-4">
              {socialLinks.map((social) => {
                const Icon = social.icon
                return (
                  <Link
                    key={social.name}
                    href={social.href}
                    className="text-muted-foreground hover:text-foreground transition-colors"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <Icon className="h-5 w-5" />
                    <span className="sr-only">{social.name}</span>
                  </Link>
                )
              })}
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}