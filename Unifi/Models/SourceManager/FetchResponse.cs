using System;

namespace Unifi.Models.SourceManager
{
    /// <summary>
    /// The response from the source
    /// </summary>
    public class FetchResponse
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string MACAddress { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public bool State { get; set; } = false;

        /// <inheritdoc />
        public override string ToString() => $"Mac: {this.MACAddress}, State: {this.State}";
    }
}