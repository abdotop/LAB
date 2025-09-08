'use client';

import { useState } from 'react';
import { UserCheck, FileText, BarChart3, Upload, Users, Eye } from 'lucide-react';
import Link from 'next/link';

export default function MandatairesPage() {
  const [user] = useState({
    name: 'Fatou Diallo',
    role: 'MANDATAIRE',
    candidate: 'Macky Sall',
    party: 'Benno Bokk Yakaar'
  });

  const stats = [
    { title: 'Parrainages Collectés', value: '1,245', icon: UserCheck, color: 'bg-green-500' },
    { title: 'Objectif', value: '2,000', icon: BarChart3, color: 'bg-blue-500' },
    { title: 'Taux de Réussite', value: '62%', icon: FileText, color: 'bg-purple-500' },
    { title: 'Régions Couvertes', value: '8/14', icon: Users, color: 'bg-orange-500' },
  ];

  const recentActivities = [
    { action: 'Parrainage validé - Dakar', time: 'Il y a 2 minutes', citizen: 'Amadou Ba' },
    { action: 'Parrainage validé - Thiès', time: 'Il y a 5 minutes', citizen: 'Awa Diop' },
    { action: 'Assistance terrain - Kaolack', time: 'Il y a 15 minutes', citizen: 'Omar Sow' },
    { action: 'Parrainage validé - Saint-Louis', time: 'Il y a 30 minutes', citizen: 'Mariama Fall' },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="w-10 h-10 bg-senegal-green rounded-lg flex items-center justify-center">
                <UserCheck className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">Espace Mandataire</h1>
                <p className="text-sm text-gray-600">Campagne de {user.candidate} - {user.party}</p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-600">Connecté en tant que {user.name}</span>
              <button className="btn-secondary">
                Déconnexion
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="p-6">
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

        {/* Main Actions */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <div className="card p-6">
            <div className="text-center">
              <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Upload className="w-8 h-8 text-blue-600" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Assistance Terrain</h3>
              <p className="text-gray-600 mb-4">
                Aidez les citoyens dans les zones à faible connectivité
              </p>
              <button className="btn-primary w-full">
                Scanner une CNI
              </button>
            </div>
          </div>

          <div className="card p-6">
            <div className="text-center">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <BarChart3 className="w-8 h-8 text-green-600" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Tableau de Bord</h3>
              <p className="text-gray-600 mb-4">
                Suivez la progression en temps réel
              </p>
              <button className="btn-secondary w-full">
                Voir les Statistiques
              </button>
            </div>
          </div>

          <div className="card p-6">
            <div className="text-center">
              <div className="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Eye className="w-8 h-8 text-purple-600" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Suivi Détaillé</h3>
              <p className="text-gray-600 mb-4">
                Consultez tous les parrainages collectés
              </p>
              <button className="btn-secondary w-full">
                Liste Complète
              </button>
            </div>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="card">
            <div className="p-6 border-b border-gray-200">
              <h3 className="text-lg font-semibold text-gray-900">Activité Récente</h3>
            </div>
            <div className="p-6">
              <div className="space-y-4">
                {recentActivities.map((activity, index) => (
                  <div key={index} className="flex items-center space-x-3">
                    <div className="w-3 h-3 bg-green-500 rounded-full" />
                    <div className="flex-1">
                      <p className="text-sm font-medium text-gray-900">{activity.action}</p>
                      <p className="text-xs text-gray-500">{activity.citizen} • {activity.time}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="card">
            <div className="p-6 border-b border-gray-200">
              <h3 className="text-lg font-semibold text-gray-900">Progression par Région</h3>
            </div>
            <div className="p-6">
              <div className="space-y-4">
                {[
                  { region: 'Dakar', count: 342, target: 500, percentage: 68 },
                  { region: 'Thiès', count: 185, target: 300, percentage: 62 },
                  { region: 'Kaolack', count: 156, target: 250, percentage: 62 },
                  { region: 'Saint-Louis', count: 123, target: 200, percentage: 62 },
                ].map((region) => (
                  <div key={region.region} className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="font-medium text-gray-900">{region.region}</span>
                      <span className="text-gray-600">{region.count}/{region.target}</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-senegal-green h-2 rounded-full transition-all duration-300"
                        style={{ width: `${region.percentage}%` }}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}