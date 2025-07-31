import React from 'react';
import { TrendingUp, AlertTriangle, Database, CheckCircle, XCircle, Clock } from 'lucide-react';
import type { Target, AlertRule } from '../types';

interface DashboardProps {
  targets: Target[];
  alertRules: AlertRule[];
}

export const Dashboard: React.FC<DashboardProps> = ({ targets, alertRules }) => {
  const stats = [
    {
      label: 'Active Targets',
      value: targets.length,
      change: '+2',
      trend: 'up',
      icon: Database,
      color: 'blue',
    },
    {
      label: 'Alert Rules',
      value: alertRules.length,
      change: '+1',
      trend: 'up',
      icon: AlertTriangle,
      color: 'yellow',
    },
    {
      label: 'Health Score',
      value: '98%',
      change: '+1%',
      trend: 'up',
      icon: TrendingUp,
      color: 'green',
    },
    {
      label: 'Last Reload',
      value: '2m ago',
      change: 'Success',
      trend: 'stable',
      icon: CheckCircle,
      color: 'green',
    },
  ];

  const recentActivity = [
    { type: 'target', action: 'Added', name: 'node-exporter', time: '5m ago', status: 'success' },
    { type: 'alert', action: 'Modified', name: alertRules[0]?.alert_name || 'HighCPUUsage', time: '12m ago', status: 'success' },
    { type: 'config', action: 'Reloaded', name: 'prometheus.yml', time: '15m ago', status: 'success' },
    { type: 'alert', action: 'Generated', name: 'DiskSpaceWarning', time: '1h ago', status: 'success' },
  ];

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
        <p className="text-gray-400">Monitor your Prometheus configuration and system health</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat, index) => {
          const Icon = stat.icon;
          const colorClasses = {
            blue: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
            yellow: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/20',
            green: 'bg-green-500/10 text-green-400 border-green-500/20',
          }[stat.color];

          return (
            <div key={index} className="bg-gray-800 rounded-xl p-6 border border-gray-700 hover:border-gray-600 transition-colors">
              <div className="flex items-center justify-between mb-4">
                <div className={`p-3 rounded-lg border ${colorClasses}`}>
                  <Icon className="w-6 h-6" />
                </div>
                <span className={`text-sm px-2 py-1 rounded-full ${
                  stat.trend === 'up' ? 'bg-green-500/10 text-green-400' : 'bg-gray-600/50 text-gray-400'
                }`}>
                  {stat.change}
                </span>
              </div>
              <div>
                <h3 className="text-2xl font-bold text-white mb-1">{stat.value}</h3>
                <p className="text-gray-400 text-sm">{stat.label}</p>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Configuration Overview */}
        <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
          <h2 className="text-xl font-semibold text-white mb-6">Configuration Overview</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-gray-700 rounded-lg">
              <div className="flex items-center gap-3">
                <Database className="w-5 h-5 text-blue-400" />
                <div>
                  <p className="text-white font-medium">Scrape Targets</p>
                  <p className="text-gray-400 text-sm">{targets.length} configured jobs</p>
                </div>
              </div>
              <CheckCircle className="w-5 h-5 text-green-400" />
            </div>
            
            <div className="flex items-center justify-between p-4 bg-gray-700 rounded-lg">
              <div className="flex items-center gap-3">
                <AlertTriangle className="w-5 h-5 text-yellow-400" />
                <div>
                  <p className="text-white font-medium">Alert Rules</p>
                  <p className="text-gray-400 text-sm">{alertRules.length} active rules</p>
                </div>
              </div>
              <CheckCircle className="w-5 h-5 text-green-400" />
            </div>

            <div className="flex items-center justify-between p-4 bg-gray-700 rounded-lg">
              <div className="flex items-center gap-3">
                <Clock className="w-5 h-5 text-purple-400" />
                <div>
                  <p className="text-white font-medium">Configuration Sync</p>
                  <p className="text-gray-400 text-sm">Auto-reload enabled</p>
                </div>
              </div>
              <CheckCircle className="w-5 h-5 text-green-400" />
            </div>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
          <h2 className="text-xl font-semibold text-white mb-6">Recent Activity</h2>
          <div className="space-y-4">
            {recentActivity.map((activity, index) => (
              <div key={index} className="flex items-center gap-4 p-4 bg-gray-700 rounded-lg">
                <div className={`p-2 rounded-lg ${
                  activity.type === 'target' ? 'bg-blue-500/20 text-blue-400' :
                  activity.type === 'alert' ? 'bg-yellow-500/20 text-yellow-400' :
                  'bg-purple-500/20 text-purple-400'
                }`}>
                  {activity.type === 'target' ? <Database className="w-4 h-4" /> :
                   activity.type === 'alert' ? <AlertTriangle className="w-4 h-4" /> :
                   <CheckCircle className="w-4 h-4" />}
                </div>
                <div className="flex-1">
                  <p className="text-white font-medium">
                    {activity.action} {activity.name}
                  </p>
                  <p className="text-gray-400 text-sm">{activity.time}</p>
                </div>
                {activity.status === 'success' ? (
                  <CheckCircle className="w-5 h-5 text-green-400" />
                ) : (
                  <XCircle className="w-5 h-5 text-red-400" />
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};