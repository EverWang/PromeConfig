import React, { useState, useEffect } from 'react';
import { AlertTriangle, Database, Settings, Zap, RefreshCw, Eye, Download, Upload, Cpu } from 'lucide-react';
import { AuthWrapper } from './components/AuthWrapper';
import { Sidebar } from './components/Sidebar';
import { Dashboard } from './components/Dashboard';
import { TargetManagement } from './components/TargetManagement';
import { AlertRuleManagement } from './components/AlertRuleManagement';
import { ConfigPreview } from './components/ConfigPreview';
import { PrometheusAPI } from './components/PrometheusAPI';
import { TargetService } from './services/targetService';
import { AlertRuleService } from './services/alertRuleService';
import type { Target, AlertRule } from './lib/supabase';

type ActiveView = 'dashboard' | 'targets' | 'alerts' | 'preview' | 'api';

function App() {
  const [activeView, setActiveView] = useState<ActiveView>('dashboard');
  const [targets, setTargets] = useState<Target[]>([]);
  const [alertRules, setAlertRules] = useState<AlertRule[]>([]);
  const [loading, setLoading] = useState(true);

  // 加载数据
  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [targetsData, alertRulesData] = await Promise.all([
        TargetService.getTargets(),
        AlertRuleService.getAlertRules()
      ]);
      setTargets(targetsData);
      setAlertRules(alertRulesData);
    } catch (error) {
      console.error('Error loading data:', error);
    } finally {
      setLoading(false);
    }
  };

  const menuItems = [
    { id: 'dashboard', label: 'Dashboard', icon: Cpu },
    { id: 'targets', label: 'Targets', icon: Database },
    { id: 'alerts', label: 'Alert Rules', icon: AlertTriangle },
    { id: 'preview', label: 'Config Preview', icon: Eye },
    { id: 'api', label: 'API Management', icon: RefreshCw },
  ];

  const renderActiveView = () => {
    if (loading) {
      return (
        <div className="flex items-center justify-center h-full">
          <div className="text-white">Loading...</div>
        </div>
      );
    }

    switch (activeView) {
      case 'dashboard':
        return <Dashboard targets={targets} alertRules={alertRules} />;
      case 'targets':
        return <TargetManagement targets={targets} onDataChange={loadData} />;
      case 'alerts':
        return <AlertRuleManagement alertRules={alertRules} onDataChange={loadData} />;
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