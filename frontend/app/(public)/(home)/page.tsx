"use client"

import { useState } from 'react'
import { Play } from 'lucide-react'
import { toast } from '@/packages/ui/use-toast'
import Image from 'next/image'

export default function LandingPage() {
  const [email, setEmail] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (email) {
      toast({
        title: 'Request Submitted!',
        status: 'success',
        description: "We'll be in touch soon.",
        duration: 4000,
      })
      setEmail('')
    } else {
      toast({
        title: 'Email Required',
        status: 'error',
        description: 'Please enter your email address.',
        duration: 4000,
      })
    }
  }

  const content = [
    {
      id: 1,
      title: 'Fewer tools',
      note: 'Reduce 10+ apps with one unified system',
      img: '/set.png'
    },
    {
      id: 2,
      title: 'Lower costs',
      note: 'Cut software spend by up to 60%',
      img: '/spend.png'
    },
    {
      id: 3,
      title: 'Faster decisions',
      note: 'Unified data means instant insights',
      img: '/fast.png'
    },
    {
      id: 4,
      title: 'Full visibility',
      note: 'See everything, everywhere, all at once',
      img: '/eye.png'
    }
  ]

  return (
    <div className="min-h-screen bg-[#0F0E11] text-white">

      {/* Header */}
      <header className="flex items-center justify-between px-6 md:px-8 py-6 max-w-7xl mx-auto">
        <Image src="/logo.png" alt="logo" width={120} height={40} />

        <button
          onClick={() => location.replace('/')}
          className="flex items-center gap-2 px-5 py-2.5 bg-black hover:bg-black/60 rounded-lg text-sm font-medium transition"
        >
          <Play className="w-4 h-4 fill-current text-black bg-amber-50 rounded-md p-0.5" />
          View Demo
        </button>
      </header>

      {/* Hero Section */}
      <main className="flex flex-col items-center text-center px-6 pt-20 pb-32">

        <h1 className="text-4xl md:text-6xl font-bold max-w-4xl leading-tight">
          One Platform to rule them all
        </h1>

        <p className="mt-6 text-gray-400 max-w-2xl">
          Alloy replaces scattered workplace tools with one intelligent operating system — chat, tasks, workflows, and payroll — all connected in one secure workplace.
        </p>

        <button className="mt-8 bg-linear-to-r from-[#4F46E5] to-[#1DA1F2] px-6 py-3 rounded-xl font-medium">
          Join the Waitlist
        </button>

        {/* Laptop Image */}
        {/* <div className="mt-16 w-full max-w-4xl">
          <Image
            src="/laptop.png"
            alt="Laptop preview"
            width={1200}
            height={700}
            className="w-full h-auto mx-auto"
          />
          <img src='/line.png' alt='underline'  className=''/>
        </div> */}

        <div className="mt-16 w-full max-w-4xl mx-auto">

  {/* Laptop (15% smaller) */}
  <div className="w-[90%] mx-auto">
    <Image
      src="/laptop.png"
      alt="Laptop preview"
      width={1600}
      height={700}
      className="w-full h-auto"
    />
  </div>

  {/* Line (Full Width) */}
  <Image
    src="/line.png"
    alt="Underline"
    width={1600}
    height={200}
    className="w-full h-auto"
  />

</div>

        {/* Features Title */}
        <h2 className="mt-24 text-2xl md:text-4xl font-bold">
          Fewer tools. Lower costs. Faster decisions
        </h2>

        <p className="mt-6 text-gray-400 max-w-2xl">
          Join the waitlist and be among the first to simplify how your business runs
        </p>

        {/* Feature Cards */}
        <div className="mt-12 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl w-full">
          {content.map((item) => (
            <div
              key={item.id}
              className="bg-[#1DA1F21A] p-6 rounded-xl text-center hover:bg-[#1DA1F230] transition"
            >
              <Image
                src={item.img}
                alt={item.title}
                width={50}
                height={50}
                className="mx-auto mb-4"
              />
              <p className="font-semibold">{item.title}</p>
              <p className="text-sm text-gray-400 mt-2">{item.note}</p>
            </div>
          ))}
        </div>

        {/* Email Form */}
        <form onSubmit={handleSubmit} className="mt-16 w-full max-w-md">
          <div className="relative">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email address"
              className="w-full bg-[#FFFFFF1A] rounded-lg border-b border-gray-600 py-3 pl-4 pr-32 text-white placeholder:text-white outline-none"
            />
            <button
              type="submit"
              className="absolute right-0 top-1/2 -translate-y-1/2 mr-1 bg-linear-to-r from-[#4F46E5] to-[#1DA1F2] px-4 py-2 rounded-lg text-sm"
            >
              Join the Waitlist
            </button>
          </div>
        </form>

      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800 py-8 px-6">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-center gap-6 text-sm text-gray-400">

          <Image src="/logo.png" alt="logo" width={100} height={30} />

          <div className="flex gap-8">
            <p className="hover:text-white cursor-pointer">Privacy</p>
            <p className="hover:text-white cursor-pointer">Terms</p>
          </div>

          <div>
            © {new Date().getFullYear()} Alloy. All rights reserved.
          </div>

        </div>
      </footer>

    </div>
  )
}