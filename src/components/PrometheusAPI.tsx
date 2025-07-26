import React, { useState } from 'react';
import { RefreshCw, Server, CheckCircle, XCircle, Settings, AlertCircle, Loader } from 'lucide-react';

export const PrometheusAPI: React.FC = () => {
  const [prometheusUrl, setPrometheusUrl] = useState('http://localhost:9090');
  const [apiKey, setApiKey] = useState('');
  const [connectionStatus, setConnectionStatus] = useState<'idle' | 'connecting' | 'connected' | 'error'>('idle');
  const [reloadStatus, setReloadStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [lastReload, setLastReload] = useState<string>('');
  const [configStatus, setConfigStatus] = useState<any>(null);

  const testConnection = async () => {
    setConnectionStatus('connecting');
    
    // Simulate API call
    setTimeout(() => {
      // Mock successful connection
      setConnectionStatus('connected');
      setConfigStatus({
        version: '2.45.0',
        uptime: '2h 15m',
        targets_active: 8,
        targets_total: 10,
        rules_loaded: 15,
        last_config_time: new Date().toISOString(),
      });
    }, 1500);
  };

  const reloadConfiguration = async () => {
    setReloadStatus('loading');
    
    // Simulate config reload
    setTimeout(() => {
      setReloadStatus('success');
      setLastReload(new Date().toLocaleString());
      
      // Reset status after 3 seconds
      setTimeout(() => setReloadStatus('idle'), 3000);
    }, 2000);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'connected':
      case 'success':
        return 'text-green-400';
      case 'error':
        return 'text-red-400';
      case 'connecting':
      case 'loading':
        return 'text-blue-400';
      default:
        return 'text-gray-400';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'connected':
      case 'success':
        return <CheckCircle className="w-5 h-5 text-green-400" />;
      case 'error':
        return <XCircle className="w-5 h-5 text-red-400" />;
      case 'connecting':
      case 'loading':
        return <Loader className="w-5 h-5 text-blue-400 animate-spin" />;
      default:
        return <Server className="w-5 h-5 text-gray-400" />;
    }
  };

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2">Prometheus API Management</h1>
        <p className="text-gray-400">Connect to Prometheus and manage configuration reloads</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Connection Configuration */}
        <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
          <h2 className="text-xl font-semibold text-white mb-6 flex items-center gap-2">
            <Server className="w-5 h-5 text-blue-400" />
            Connection Configuration
          </h2>
          
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Prometheus URL
              </label>
              <input
                type="url"
                value={prometheusUrl}
                onChange={(e) => setPrometheusUrl(e.target.value)}
                className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="http://localhost:9090"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                API Key (Optional)
              </label>
              <input
                type="password"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                className="w-full px-4 py-3 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter API key if required"
              />
            </div>

            <div className="flex items-center justify-between p-4 bg-gray-700 rounded-lg">
              <div className="flex items-center gap-3">
                {getStatusIcon(connectionStatus)}
                <div>
                  <p className="text-white font-medium">Connection Status</p>
                  <p className={`text-sm capitalize ${getStatusColor(connectionStatus)}`}>
                    {connectionStatus === 'idle' ? 'Not connected' : connectionStatus}
                  </p>
                </div>
              </div>
              <button
                onClick={testConnection}
                disabled={connectionStatus === 'connecting'}
                className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition-colors"
              >
                {connectionStatus === 'connecting' ? (
                  <Loader className="w-4 h-4 animate-spin" />
                ) : (
                  <RefreshCw className="w-4 h-4" />
                )}
                Test Connection
              </button>
            </div>
          </div>
        </div>

        {/* Configuration Management */}
        <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
          <h2 className="text-xl font-semibold text-white mb-6 flex items-center gap-2">
            <Settings className="w-5 h-5 text-yellow-400" />
            Configuration Management
          </h2>

          <div className="space-y-6">
            <div className="bg-gray-700 rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-white">Configuration Reload</h3>
                {getStatusIcon(reloadStatus)}
              </div>
              
              <p className="text-gray-300 text-sm mb-4">
                Reload the Prometheus configuration without restarting the service. This will apply 
                any changes made to targets and alert rules.
              </p>

              <button
                onClick={reloadConfiguration}
                disabled={reloadStatus === 'loading' || connectionStatus !== 'connected'}
                className="w-full bg-yellow-600 hover:bg-yellow-700 disabled:bg-gray-600 text-white px-4 py-3 rounded-lg flex items-center justify-center gap-2 transition-colors"
              >
                {reloadStatus === 'loading' ? (
                  <>
                    <Loader className="w-4 h-4 animate-spin" />
                    Reloading Configuration...
                  </>
                ) : (
                  <>
                    <RefreshCw className="w-4 h-4" />
                    Reload Configuration
                  </>
                )}
              </button>

              {lastReload && (
                <p className="text-gray-400 text-sm mt-3">
                  Last reload: {lastReload}
                </p>
              )}
            </div>

            {/* Health Check */}
            <div className="bg-gray-700 rounded-lg p-4">
              <h3 className="text-lg font-medium text-white mb-3 flex items-center gap-2">
                <AlertCircle className="w-5 h-5 text-blue-400" />
                Health Check
              </h3>
              
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-gray-400">Status</p>
                  <p className="text-green-400 font-medium">Healthy</p>
                </div>
                <div>
                  <p className="text-gray-400">Uptime</p>
                  <p className="text-white font-medium">2h 15m</p>
                </div>
                <div>
                  <p className="text-gray-400">Active Targets</p>
                  <p className="text-white font-medium">8/10</p>
                </div>
                <div>
                  <p className="text-gray-400">Rules Loaded</p>
                  <p className="text-white font-medium">15</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Prometheus Status Overview */}
      {configStatus && (
        <div className="mt-8 bg-gray-800 rounded-xl p-6 border border-gray-700">
          <h2 className="text-xl font-semibold text-white mb-6">Prometheus Status Overview</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-500/20 rounded-lg">
                  <Server className="w-5 h-5 text-blue-400" />
                </div>
                <div>
                  <p className="text-gray-400 text-sm">Version</p>
                  <p className="text-white font-semibold">{configStatus.version}</p>
                </div>
              </div>
            </div>

            <div className="bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-green-500/20 rounded-lg">
                  <CheckCircle className="w-5 h-5 text-green-400" />
                </div>
                <div>
                  <p className="text-gray-400 text-sm">Uptime</p>
                  <p className="text-white font-semibold">{configStatus.uptime}</p>
                </div>
              </div>
            </div>

            <div className="bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-yellow-500/20 rounded-lg">
                  <AlertTriangle className="w-5 h-5 text-yellow-400" />
                </div>
                <div>
                  <p className="text-gray-400 text-sm">Active Targets</p>
                  <p className="text-white font-semibold">
                    {configStatus.targets_active}/{configStatus.targets_total}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-gray-700 rounded-lg p-4">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-purple-500/20 rounded-lg">
                  <Settings className="w-5 h-5 text-purple-400" />
                </div>
                <div>
                  <p className="text-gray-400 text-sm">Rules Loaded</p>
                  <p className="text-white font-semibold">{configStatus.rules_loaded}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* API Endpoints */}
      <div className="mt-8 bg-gray-800 rounded-xl p-6 border border-gray-700">
        <h2 className="text-xl font-semibold text-white mb-6">Available API Endpoints</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-white font-medium mb-2">Configuration Reload</h3>
            <code className="text-blue-300 text-sm">POST /-/reload</code>
            <p className="text-gray-400 text-sm mt-2">Reload the configuration file</p>
          </div>

          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-white font-medium mb-2">Health Check</h3>
            <code className="text-blue-300 text-sm">GET /-/healthy</code>
            <p className="text-gray-400 text-sm mt-2">Check Prometheus health status</p>
          </div>

          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-white font-medium mb-2">Ready Check</h3>
            <code className="text-blue-300 text-sm">GET /-/ready</code>
            <p className="text-gray-400 text-sm mt-2">Check if Prometheus is ready</p>
          </div>

          <div className="bg-gray-700 rounded-lg p-4">
            <h3 className="text-white font-medium mb-2">Configuration</h3>
            <code className="text-blue-300 text-sm">GET /api/v1/status/config</code>
            <p className="text-gray-400 text-sm mt-2">Get current configuration</p>
          </div>
        </div>
      </div>
    </div>
  );
};