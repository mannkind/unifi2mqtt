using System;

namespace Unifi.Models.Source
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public record Response
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string MACAddress { get; init; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool State { get; init; } = false;
    }
}