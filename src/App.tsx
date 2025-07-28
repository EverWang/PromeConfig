import React, { useState } from 'react';
import { AlertTriangle, Database, Settings, Zap, RefreshCw, Eye, Download, Upload, Cpu } from 'lucide-react';
import { AuthWrapper } from './components/AuthWrapper';
import { Sidebar } from './components/Sidebar';
import { Dashboard } from './components/Dashboard';
import { TargetManagement } from './components/TargetManagement';
import { AlertRuleManagement } from './components/AlertRuleManagement';
import { ConfigPreview } from './components/ConfigPreview';
import { PrometheusAPI } from './components/PrometheusAPI';

type ActiveView = 'dashboard' | 'targets' | 'alerts' | 'preview' | 'api';

function App() {
  const [activeView, setActiveView] = useState<ActiveView>('dashboard');
  const [targets, setTargets] = useState([
    {
      id: '1',
      job_name: 'node-exporter',
      static_configs: [{ targets: ['localhost:9100'] }],
      scrape_interval: '15s',
      metrics_path: '/metrics',
    },
    {
      id: '2',
      job_name: 'prometheus',
      static_configs: [{ targets: ['localhost:9090'] }],
      scrape_interval: '15s',
      metrics_path: '/metrics',
    },
  ]);

  const [alertRules, setAlertRules] = useState([
    {
      id: '1',
      alert: 'HighCPUUsage',
      expr: 'cpu_usage_percent > 80',
      for: '5m',
      labels: { severity: 'warning' },
      annotations: {
        summary: 'High CPU usage detected',
        description: 'CPU usage is above 80% for more than 5 minutes',
      },
    },
    {
      id: '2',
      alert: 'ServiceDown',
      expr: 'up == 0',
      for: '1m',
      labels: { severity: 'critical' },
      annotations: {
        summary: 'Service is down',
        description: 'Service {{ $labels.instance }} is down',
      },
    },
  ]);

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: Cpu },
    { id: 'targets', label: 'Targets', icon: Database },
    { id: 'alerts', label: 'Alert Rules', icon: AlertTriangle },
    { id: 'preview', label: 'Config Preview', icon: Eye },
    { id: 'api', label: 'API Management', icon: RefreshCw },
  ];

  const renderActiveView = () => {
    switch (activeView) {
      case 'dashboard':
        return <Dashboard targets={targets} alertRules={alertRules} />;
      case 'targets':
        return <TargetManagement targets={targets} setTargets={setTargets} />;
      case 'alerts':
        return <AlertRuleManagement alertRules={alertRules} setAlertRules={setAlertRules} />;
      case 'preview':
        return <ConfigPreview targets={targets} alertRules={alertRules} />;
      case 'api':
        return <PrometheusAPI />;
      default:
        return <Dashboard targets={targets} alertRules={alertRules} />;
    }
  };

  return (
    <AuthWrapper>
      <div className="flex h-screen bg-gray-900 text-gray-100">
        <Sidebar
          menuItems={menuItems}
          activeView={activeView}
          setActiveView={setActiveView}
        />
        <main className="flex-1 overflow-auto">
          {renderActiveView()}
        </main>
      </div>
    </AuthWrapper>
  );
}

export default App;