'use client';

import { useState } from 'react';
import { Shield, Users, FileText, BarChart3, Settings, LogOut } from 'lucide-react';
import Link from 'next/link';

export default function DashboardPage() {
  const [user] = useState({
    name: 'Admin DGE',
    role: 'SUPER_ADMIN_DGE',
    email: 'admin@dge.sn'
  });

  const stats = [
    { title: 'Cycles Électoraux Actifs', value: '2', icon: FileText, color: 'bg-blue-500' },
    { title: 'Candidatures en Attente', value: '15', icon: Users, color: 'bg-yellow-500' },
    { title: 'Parrainages Validés', value: '12,847', icon: Shield, color: 'bg-green-500' },
    { title: 'Signalements', value: '3', icon: BarChart3, color: 'bg-red-500' },
  ];

  const menuItems = [
    { title: 'Tableau de Bord', href: '/', icon: BarChart3, active: true },
    { title: 'Cycles Électoraux', href: '/cycles', icon: FileText },
    { title: 'Candidatures', href: '/candidatures', icon: Users },
    { title: 'Parrainages', href: '/parrainages', icon: Shield },
    { title: 'Audit Trail', href: '/audit', icon: Settings },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Sidebar */}
      <div className="fixed inset-y-0 left-0 w-64 bg-white shadow-lg">
        <div className="flex h-16 items-center justify-center border-b border-gray-200">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-senegal-green rounded-lg flex items-center justify-center">
              <Shield className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold text-gray-900">DGE Admin</span>
          </div>
        </div>

        <nav className="mt-8 px-4">
          <ul className="space-y-2">
            {menuItems.map((item) => (
              <li key={item.href}>
                <Link
                  href={item.href}
                  className={`flex items-center space-x-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                    item.active
                      ? 'bg-senegal-green text-white'
                      : 'text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  <item.icon className="w-5 h-5" />
                  <span>{item.title}</span>
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <div className="absolute bottom-4 left-4 right-4">
          <div className="bg-gray-100 rounded-lg p-3">
            <div className="text-sm font-medium text-gray-900">{user.name}</div>
            <div className="text-xs text-gray-500">{user.email}</div>
            <button className="mt-2 flex items-center space-x-2 text-sm text-gray-700 hover:text-gray-900">
              <LogOut className="w-4 h-4" />
              <span>Déconnexion</span>
            </button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="ml-64">
        {/* Header */}
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="px-6 py-4">
            <div className="flex items-center justify-between">
              <h1 className="text-2xl font-bold text-gray-900">Tableau de Bord</h1>
              <div className="flex items-center space-x-4">
                <div className="text-sm text-gray-500">
                  Direction Générale des Élections - République du Sénégal
                </div>
              </div>
            </div>
          </div>
        </header>

        {/* Dashboard Content */}
        <main className="p-6">
          {/* Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            {stats.map((stat) => (
              <div key={stat.title} className="card p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">{stat.title}</p>
                    <p className="text-3xl font-bold text-gray-900">{stat.value}</p>
                  </div>
                  <div className={`w-12 h-12 ${stat.color} rounded-lg flex items-center justify-center`}>
                    <stat.icon className="w-6 h-6 text-white" />
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Recent Activity */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="card">
              <div className="p-6 border-b border-gray-200">
                <h3 className="text-lg font-semibold text-gray-900">Activités Récentes</h3>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  {[
                    { action: 'Nouvelle candidature déposée', time: 'Il y a 5 minutes', status: 'pending' },
                    { action: 'Candidature approuvée', time: 'Il y a 15 minutes', status: 'success' },
                    { action: 'Signalement de fraude traité', time: 'Il y a 1 heure', status: 'warning' },
                  ].map((activity, index) => (
                    <div key={index} className="flex items-center space-x-3">
                      <div className={`w-3 h-3 rounded-full ${
                        activity.status === 'pending' ? 'bg-yellow-500' :
                        activity.status === 'success' ? 'bg-green-500' : 'bg-orange-500'
                      }`} />
                      <div className="flex-1">
                        <p className="text-sm font-medium text-gray-900">{activity.action}</p>
                        <p className="text-xs text-gray-500">{activity.time}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="card">
              <div className="p-6 border-b border-gray-200">
                <h3 className="text-lg font-semibold text-gray-900">Actions Rapides</h3>
              </div>
              <div className="p-6">
                <div className="space-y-3">
                  <button className="btn-primary w-full py-3">
                    Nouveau Cycle Électoral
                  </button>
                  <button className="btn-secondary w-full py-3">
                    Exporter les Données
                  </button>
                  <button className="btn-secondary w-full py-3">
                    Consulter l'Audit Trail
                  </button>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}