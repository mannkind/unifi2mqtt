using System;
using System.Net.Http;
using System.Threading;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Options;

namespace Unifi
{
    /// <summary>
    /// Extensions for classes implementing IServiceCollection
    /// </summary>
    public static class IServiceCollectionExt
    {
        /// <summary>
        /// Hopefully this exists only until KoenZomers.UniFi.Api is updated.
        /// </summary>
        /// <param name="services"></param>
        /// <typeparam name="T"></typeparam>
        /// <returns></returns>
        public static IServiceCollection AddTypeNamedHttpClient<T>(this IServiceCollection services, bool allowAutoRedirect = true, TimeSpan? lifetime = null)
            where T : class =>
            services
                .AddHttpClient(typeof(T).Name)
                .ConfigurePrimaryHttpMessageHandler((x) =>
                {
                    var opts = x.GetRequiredService<IOptions<Models.Options.SourceOpts>>();
                    return SetupHttpClientHandler(allowAutoRedirect, !opts.Value.DisableSslValidation);
                })
                .SetHandlerLifetime(lifetime ?? TimeSpan.FromMinutes(2))
                .Services;

        public static HttpClientHandler SetupHttpClientHandler(bool allowAutoRedirect = true, bool validateSSL = false)
        {
            var handler = new HttpClientHandler
            {
                AllowAutoRedirect = allowAutoRedirect
            };

            if (!validateSSL)
            {
                handler.ClientCertificateOptions = ClientCertificateOption.Manual;
                handler.ServerCertificateCustomValidationCallback = (httpRequestMessage, cert, cetChain, policyErrors) => true;
            }

            return handler;
        }
    }
}